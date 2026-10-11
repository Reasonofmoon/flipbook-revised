"""Build self-hosted, unicode-range-split web fonts for Achmage SameTone.

The upstream web zips ship one ~1.4 MB WOFF2 per weight (all 11,172 Hangul
syllables). Browsers would download the whole file for the first Korean
character, so each weight is split into ranges and only the ranges a page
uses are fetched.

Usage (needs `pip install fonttools brotli`):
  python scripts/fonts/build_sametone.py <Text-woff2.zip> <Display-woff2.zip> web/static/fonts/sametone

Source: https://achmage-sametone.achmage2.chatgpt.site/ (RC2.1, SIL OFL 1.1).
License and notice files from the zip are copied next to the fonts, as the
OFL requires.
"""
import io
import pathlib
import posixpath
import sys
import zipfile

from fontTools import subset
from fontTools.ttLib import TTFont

# Weights actually used by the site's CSS.
WEIGHTS = {"Text": [400, 500, 600, 700, 800], "Display": [700]}
FAMILY = {"Text": "Achmage SameTone Text", "Display": "Achmage SameTone Display"}

# The 2,350 Hangul syllables of KS X 1001 cover almost all everyday Korean.
COMMON_HANGUL = sorted({ord(bytes([hi, lo]).decode("euc-kr"))
                        for hi in range(0xB0, 0xC9) for lo in range(0xA1, 0xFF)})
ALL_HANGUL = range(0xAC00, 0xD7A4)

# Order matters: ko-rest is declared before ko-common. Both get the whole
# Hangul block as their CSS unicode-range, and browsers check overlapping
# faces last-declared first, so a syllable is looked up in the small common
# file and only falls through to ko-rest when it is missing there. This keeps
# the CSS tiny instead of listing thousands of scattered code points.
RANGES = {
    "latin": list(range(0x20, 0x7F)) + list(range(0xA0, 0x180)) + list(range(0x2000, 0x2070))
    + [0x20A9, 0x20AC, 0x2122, 0x2190, 0x2191, 0x2192, 0x2193, 0x2212, 0x2713, 0x00B7],
    "ko-misc": list(range(0x1100, 0x1200)) + list(range(0x3000, 0x3040)) + list(range(0x3130, 0x3190))
    + list(range(0xFF00, 0xFFF0)) + list(range(0x2460, 0x2500)),  # jamo, CJK punct, fullwidth, circled
    "ko-rest": [c for c in ALL_HANGUL if c not in set(COMMON_HANGUL)],
    "ko-common": COMMON_HANGUL,
}
CSS_RANGE_OVERRIDE = {"ko-rest": "U+AC00-D7A3", "ko-common": "U+AC00-D7A3"}

LICENSE_FILES = ("OFL.txt", "NOTICE.txt", "LICENSE", "FONTLOG.md")


def unicode_range(codepoints):
    """Compress code points into a CSS unicode-range list."""
    cps = sorted(set(codepoints))
    parts, start, prev = [], cps[0], cps[0]
    for c in cps[1:] + [None]:
        if c is not None and c == prev + 1:
            prev = c
            continue
        parts.append(f"U+{start:X}" if start == prev else f"U+{start:X}-{prev:X}")
        if c is not None:
            start = prev = c
    return ", ".join(parts)


def member(zf, suffix):
    matches = [n for n in zf.namelist() if n.endswith(suffix) and ".." not in posixpath.normpath(n).split("/")]
    if len(matches) != 1:
        sys.exit(f"expected one member ending in {suffix!r}, found {matches}")
    return zf.read(matches[0])


def main(text_zip, display_zip, out_dir):
    out = pathlib.Path(out_dir)
    (out / "licenses").mkdir(parents=True, exist_ok=True)
    css = ["/* Achmage SameTone RC2.1 — SIL Open Font License 1.1, see OFL.txt and NOTICE.txt.",
           "   Built by scripts/fonts/build_sametone.py (unicode-range split). */"]
    total = 0
    for style, zpath in (("Text", text_zip), ("Display", display_zip)):
        with zipfile.ZipFile(zpath) as zf:
            if style == "Text":  # license files are identical in both zips
                for name in LICENSE_FILES:  # top-level files, matched exactly
                    (out / name).write_bytes(zf.read(name))
                for n in zf.namelist():
                    if n.startswith("licenses/") and n.endswith(".txt"):
                        (out / "licenses" / posixpath.basename(n)).write_bytes(zf.read(n))
            for w in WEIGHTS[style]:
                src = member(zf, f"{style}-{w}.woff2")  # fonts/woff2/AchmageSameToneText-400.woff2
                for rname, cps in RANGES.items():
                    font = TTFont(io.BytesIO(src))
                    available = set(font.getBestCmap())
                    keep = [c for c in cps if c in available]
                    if not keep:
                        continue
                    opts = subset.Options()
                    opts.flavor = "woff2"
                    opts.layout_features = ["*"]
                    opts.name_IDs = ["*"]  # keep copyright/license name records
                    opts.notdef_outline = True
                    sub = subset.Subsetter(opts)
                    sub.populate(unicodes=keep)
                    sub.subset(font)
                    fname = f"SameTone{style}-{w}-{rname}.woff2"
                    font.save(out / fname)
                    size = (out / fname).stat().st_size
                    total += size
                    css.append(
                        f"@font-face {{ font-family: \"{FAMILY[style]}\"; font-style: normal; font-weight: {w}; "
                        f"font-display: swap; src: url(\"{fname}\") format(\"woff2\"); "
                        f"unicode-range: {CSS_RANGE_OVERRIDE.get(rname) or unicode_range(keep)}; }}")
                    print(f"{fname:40s} {size // 1024:5d} KB  ({len(keep)} chars)")
    (out / "sametone.css").write_text("\n".join(css) + "\n", encoding="utf-8")
    print(f"total {total // 1024} KB -> {out / 'sametone.css'}")


if __name__ == "__main__":
    if len(sys.argv) != 4:
        sys.exit(__doc__)
    main(*sys.argv[1:])
