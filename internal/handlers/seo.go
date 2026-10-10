package handlers

import (
	"unicode"
	"unicode/utf8"
)

// koreanLetterRatio is the share of Hangul among all letters at which a
// flipbook is tagged as Korean. Korean-medium English lessons mix glosses
// like "abundant 풍부한", so a minority of Hangul still marks the audience.
const koreanLetterRatio = 0.2

// detectLang returns the BCP 47 language tag for the <html lang> attribute,
// based on the title and extracted page text.
func detectLang(title string, pageTexts []string) string {
	var letters, hangul int
	count := func(s string) {
		for _, r := range s {
			if !unicode.IsLetter(r) {
				continue
			}
			letters++
			if unicode.Is(unicode.Hangul, r) {
				hangul++
			}
		}
	}
	count(title)
	for _, t := range pageTexts {
		count(t)
	}
	if letters > 0 && float64(hangul)/float64(letters) >= koreanLetterRatio {
		return "ko"
	}
	return "en"
}

// truncateRunes shortens s to at most n runes so multi-byte text (e.g. Hangul)
// is never cut in the middle of a character. It reports whether s was cut.
func truncateRunes(s string, n int) (string, bool) {
	if utf8.RuneCountInString(s) <= n {
		return s, false
	}
	return string([]rune(s)[:n]), true
}
