"""Record real-product GIFs of the local flipbook (localhost:8080) for the landing page.

Requires the docker compose stack running, a converted sample flipbook, and
`pip install playwright pillow` (uses the installed Google Chrome).

Usage: python record_gifs.py <out_dir> <slug> <sample_pptx> [names...]
  e.g. python scripts/landing/record_gifs.py web/static/img/landing unit-3-the-honeybee-dance sample.pptx
The admin password is read from the local container env inside this script
and never printed.
"""
import pathlib
import subprocess
import sys
import time

from PIL import Image
from playwright.sync_api import sync_playwright

BASE = "http://localhost:8080"
OUT = pathlib.Path(sys.argv[1])
SLUG = sys.argv[2]
PPTX = sys.argv[3]
ONLY = set(sys.argv[4:])
REPO = str(pathlib.Path(__file__).resolve().parents[2])  # repo root, for docker compose


class Recorder:
    def __init__(self, page, clip=None):
        self.page, self.clip, self.frames = page, clip, []

    def shoot(self, seconds, interval=0.08):
        end = time.time() + seconds
        while time.time() < end:
            t = time.time()
            png = self.page.screenshot(clip=self.clip, animations="allow")
            self.frames.append((png, interval))
            time.sleep(max(0, interval - (time.time() - t)))

    def hold(self, seconds):
        """Freeze the current frame for a while (cheap pause)."""
        png = self.page.screenshot(clip=self.clip)
        self.frames.append((png, seconds))

    def save(self, name, width, colors=96):
        import io
        imgs, durs = [], []
        for png, d in self.frames:
            im = Image.open(io.BytesIO(png)).convert("RGB")
            if im.width != width:
                im = im.resize((width, round(im.height * width / im.width)), Image.LANCZOS)
            imgs.append(im)
            durs.append(int(d * 1000))
        palette_src = imgs[len(imgs) // 2].quantize(colors=colors, method=Image.Quantize.MEDIANCUT)
        frames = [im.quantize(palette=palette_src, dither=Image.Dither.NONE) for im in imgs]
        path = OUT / f"{name}.gif"
        frames[0].save(path, save_all=True, append_images=frames[1:], duration=durs,
                       loop=0, optimize=True, disposal=1)
        print(f"{name}.gif  {len(frames)} frames  {path.stat().st_size / 1024:.0f} KB  {frames[0].size}")


def admin_password():
    return subprocess.run(
        ["docker", "compose", "exec", "-T", "app", "sh", "-c", 'printf %s "$FLIPBOOK_ADMIN_PASSWORD"'],
        cwd=REPO, capture_output=True, text=True, check=True).stdout


def drag_flip(page, rec, box, forward=True, steps=9):
    """Flip one page by dragging the book's corner, capturing the curl."""
    x, y, w, h = box["x"], box["y"], box["width"], box["height"]
    sx, ex = (x + w - 6, x - w * 0.25) if forward else (x + 6, x + w * 1.25)
    sy, ey = y + h * 0.86, y + h * 0.62
    page.mouse.move(sx, sy)
    page.mouse.down()
    for i in range(1, steps + 1):
        t = i / steps
        page.mouse.move(sx + (ex - sx) * t, sy + (ey - sy) * t)
        rec.shoot(0.085, 0.085)
    page.mouse.up()
    rec.shoot(0.45, 0.09)


def book_box(page):
    return page.locator(".stf__parent").bounding_box()


def want(name):
    return not ONLY or name in ONLY


def wait_viewer(page):
    page.wait_for_selector(".stf__parent", timeout=15000)
    page.wait_for_timeout(1200)


def hero_flip(browser):
    ctx = browser.new_context(viewport={"width": 1280, "height": 760}, locale="ko-KR")
    page = ctx.new_page()
    page.goto(f"{BASE}/v/{SLUG}")
    wait_viewer(page)
    rec = Recorder(page)
    rec.hold(1.4)
    box = book_box(page)
    for _ in range(5):
        drag_flip(page, rec, book_box(page) or box)
        rec.hold(1.1)
    page.keyboard.press("Home")  # jump back to the cover so the loop restarts cleanly
    page.evaluate("window.scrollTo(0,0)")
    for _ in range(6):
        page.keyboard.press("ArrowLeft")
    page.wait_for_timeout(400)
    rec.shoot(0.3, 0.1)
    rec.save("hero-flip", 880, colors=80)
    ctx.close()


def step_upload(browser):
    ctx = browser.new_context(viewport={"width": 1100, "height": 720}, locale="ko-KR")
    page = ctx.new_page()
    page.goto(f"{BASE}/login")
    page.fill("#password", admin_password())
    page.click("button[type=submit]")
    page.wait_for_url(f"{BASE}/admin")
    page.goto(f"{BASE}/admin/upload")
    page.wait_for_timeout(600)
    rec = Recorder(page)
    rec.hold(1.0)
    page.fill("#title", "")
    page.type("#title", "Unit 3. The Honeybee Dance", delay=45)
    rec.shoot(0.4)
    page.set_input_files("#file", PPTX)
    rec.shoot(0.6)
    rec.hold(0.7)
    page.click("#submit-btn")
    # Conversion takes ~10-20s; sample sparsely, then show the finished state.
    end = time.time() + 150
    while time.time() < end:
        png = page.screenshot()
        rec.frames.append((png, 0.25))  # time-lapse: 1s of waiting shown as 0.25s
        if page.locator("#step-done.completed").count():
            break
        time.sleep(1.0)
    page.wait_for_timeout(400)
    rec.hold(2.6)
    rec.save("step-upload", 880, colors=64)
    ctx.close()


def step_share(browser):
    ctx = browser.new_context(viewport={"width": 1100, "height": 720}, locale="ko-KR")
    page = ctx.new_page()
    page.goto(f"{BASE}/v/{SLUG}")
    wait_viewer(page)
    rec = Recorder(page)
    rec.hold(0.9)
    page.click("#btn-share")
    rec.shoot(0.6)
    rec.hold(1.2)
    page.click("#share-modal .copy-field button")
    rec.shoot(0.5)
    rec.hold(1.0)
    page.locator("#share-modal .copy-field button").nth(1).click()
    rec.shoot(0.5)
    rec.hold(1.6)
    rec.save("step-share", 880, colors=64)
    ctx.close()


def step_mobile(browser, pw):
    device = pw.devices["iPhone 13 landscape"]
    ctx = browser.new_context(**device, locale="ko-KR")
    page = ctx.new_page()
    page.goto(f"{BASE}/v/{SLUG}")
    wait_viewer(page)
    rec = Recorder(page)
    rec.hold(1.0)
    # Flip gestures are disabled on touch devices by design; students tap the arrows.
    for _ in range(4):
        page.tap("#nav-next")
        rec.shoot(0.35, 0.09)
        rec.hold(1.0)
    rec.save("step-mobile", 640, colors=96)
    ctx.close()


def qa_search(browser):
    ctx = browser.new_context(viewport={"width": 1100, "height": 700}, locale="ko-KR")
    page = ctx.new_page()
    page.goto(f"{BASE}/v/{SLUG}")
    wait_viewer(page)
    rec = Recorder(page)
    rec.hold(0.9)
    page.click("#btn-search")
    rec.shoot(0.4)
    for ch in "꽃까지의거리":
        page.keyboard.insert_text(ch)
        page.dispatch_event("#search-input", "input")
        rec.shoot(0.25)
    rec.shoot(0.9, 0.07)
    rec.hold(2.2)
    rec.save("qa-search", 880, colors=96)
    ctx.close()


BLOG_HTML = """<!doctype html><html lang="ko"><head><meta charset="utf-8">
<style>
 body{margin:0;background:#f3f4f6;font-family:'Malgun Gothic',sans-serif;color:#222}
 .bar{background:#fff;border-bottom:1px solid #e5e7eb;padding:14px 28px;font-weight:700}
 .post{max-width:760px;margin:24px auto;background:#fff;padding:28px 32px;border-radius:6px}
 h1{font-size:24px;margin:0 0 6px} .meta{color:#888;font-size:13px;margin-bottom:18px}
 p{line-height:1.7;font-size:15px}
</style></head><body>
<div class="bar">우리 학원 영어 블로그</div>
<div class="post"><h1>[중2] Unit 3 꿀벌의 춤 — 수업 자료</h1>
<div class="meta">2026. 10. 11 · 수업 자료</div>
<p>오늘 수업에서 다룬 교재입니다. 아래에서 바로 넘겨 보며 복습하세요.</p>
<iframe src="__SRC__" width="696" height="430" frameborder="0" allowfullscreen style="border:none;border-radius:8px"></iframe>
<p>다음 시간에는 Unit 4를 함께 읽어요.</p></div></body></html>"""


def qa_embed(browser):
    html = BLOG_HTML.replace("__SRC__", f"{BASE}/embed/{SLUG}")
    ctx = browser.new_context(viewport={"width": 1000, "height": 640}, locale="ko-KR")
    page = ctx.new_page()
    page.route(f"{BASE}/__blog-demo", lambda r: r.fulfill(body=html, content_type="text/html; charset=utf-8"))
    page.goto(f"{BASE}/__blog-demo")
    frame = page.frame_locator("iframe")
    frame.locator(".stf__parent").wait_for(timeout=15000)
    page.wait_for_timeout(1200)
    rec = Recorder(page)
    rec.hold(1.2)
    for _ in range(3):
        drag_flip(page, rec, frame.locator(".stf__parent").bounding_box(), steps=9)
        rec.hold(0.9)
    rec.save("qa-embed", 800, colors=96)
    ctx.close()


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    with sync_playwright() as pw:
        browser = pw.chromium.launch(channel="chrome", headless=True)
        if want("hero-flip"):
            hero_flip(browser)
        if want("step-upload"):
            step_upload(browser)
        if want("step-share"):
            step_share(browser)
        if want("step-mobile"):
            step_mobile(browser, pw)
        if want("qa-search"):
            qa_search(browser)
        if want("qa-embed"):
            qa_embed(browser)
        browser.close()


if __name__ == "__main__":
    main()
