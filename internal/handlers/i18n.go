package handlers

// uiStrings holds the viewer's user-facing text for one language.
type uiStrings map[string]string

// uiText maps a language tag (as returned by detectLang) to viewer UI text.
// Every language must define every key in "en"; see TestUITextComplete.
var uiText = map[string]uiStrings{
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
	},
}

// uiFor returns the UI text for lang, falling back to English.
func uiFor(lang string) uiStrings {
	if t, ok := uiText[lang]; ok {
		return t
	}
	return uiText["en"]
}
