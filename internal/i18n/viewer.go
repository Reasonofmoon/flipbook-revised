// Package i18n holds user-facing UI text for the viewer and admin pages.
package i18n

// Strings holds user-facing text for one language, keyed by message ID.
type Strings map[string]string

// viewerText maps a language tag to viewer UI text.
// Every language must define every key in "en"; see TestTablesComplete.
var viewerText = map[string]Strings{
	"en": {
		"prev_page":          "Previous page",
		"next_page":          "Next page",
		"prev_page_hint":     "Previous page (Left Arrow)",
		"next_page_hint":     "Next page (Right Arrow)",
		"search_hint":        "Search (Ctrl+F)",
		"grid_hint":          "Grid view (G)",
		"fullscreen_hint":    "Fullscreen (F)",
		"share_hint":         "Share / Embed",
		"created_before":     "Created with ",
		"created_after":      "",
		"share_title":        "Share this flipbook",
		"share_from_page":    "Share from current page",
		"link":               "Link",
		"copy":               "Copy",
		"copied":             "Copied!",
		"embed_code":         "Embed code",
		"close":              "Close",
		"search_placeholder": "Search in slides...",
		"prev_match":         "Previous match",
		"next_match":         "Next match",
		"close_search":       "Close search",
		"no_results":         "No results",
		"all_pages":          "All Pages",
		"close_grid":         "Close grid",
		"slide":              "Slide",
		"page":               "Page",
		"loading":            "Loading...",
		"preparing":          "This flipbook is being prepared. The page will refresh automatically.",
		"rotate_hint":        "Rotate for best viewing experience",
	},
	"ko": {
		"prev_page":          "이전 페이지",
		"next_page":          "다음 페이지",
		"prev_page_hint":     "이전 페이지 (←)",
		"next_page_hint":     "다음 페이지 (→)",
		"search_hint":        "검색 (Ctrl+F)",
		"grid_hint":          "전체 페이지 보기 (G)",
		"fullscreen_hint":    "전체 화면 (F)",
		"share_hint":         "공유 / 임베드",
		"created_before":     "",
		"created_after":      "으로 제작",
		"share_title":        "플립북 공유하기",
		"share_from_page":    "현재 페이지부터 공유",
		"link":               "링크",
		"copy":               "복사",
		"copied":             "복사됨!",
		"embed_code":         "임베드 코드",
		"close":              "닫기",
		"search_placeholder": "슬라이드에서 검색...",
		"prev_match":         "이전 결과",
		"next_match":         "다음 결과",
		"close_search":       "검색 닫기",
		"no_results":         "검색 결과 없음",
		"all_pages":          "전체 페이지",
		"close_grid":         "전체 페이지 닫기",
		"slide":              "슬라이드",
		"page":               "페이지",
		"loading":            "준비 중...",
		"preparing":          "플립북을 준비하고 있습니다. 완료되면 페이지가 자동으로 새로고침됩니다.",
		"rotate_hint":        "가로로 돌리면 더 크게 볼 수 있어요",
	},
}

// ResolveViewer picks the viewer UI language: the server-wide setting when it
// names a language ("ko", "en"), otherwise ("auto") the content language.
func ResolveViewer(setting, contentLang string) string {
	if _, ok := viewerText[setting]; ok {
		return setting
	}
	return contentLang
}

// Viewer returns the viewer UI text for lang, falling back to English.
func Viewer(lang string) Strings {
	return lookup(viewerText, lang)
}

func lookup(table map[string]Strings, lang string) Strings {
	if t, ok := table[lang]; ok {
		return t
	}
	return table["en"]
}
