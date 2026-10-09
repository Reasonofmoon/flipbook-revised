package handlers

import (
	"strings"
	"testing"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"english unchanged", "My Deck 2024!", "my-deck-2024"},
		{"apostrophe dropped", "Don't Panic", "dont-panic"},
		{"korean kept", "MCP 한글 테스트", "mcp-한글-테스트"},
		{"korean only", "수능 영어 독해", "수능-영어-독해"},
		{"korean punctuation", "2장: 핵심 어휘 (정리)", "2장-핵심-어휘-정리"},
		{"nfd normalized to nfc", norm.NFD.String("한글 교재"), "한글-교재"},
		{"symbols only falls back", "!!! ???", "flipbook"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slugify(tt.in); got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSlugifyTruncatesOnRuneBoundary(t *testing.T) {
	got := slugify(strings.Repeat("가나 ", 50))
	if n := utf8.RuneCountInString(got); n > maxSlugRunes {
		t.Fatalf("slug has %d runes, want <= %d", n, maxSlugRunes)
	}
	if !utf8.ValidString(got) {
		t.Fatal("slug is not valid UTF-8")
	}
	if strings.HasSuffix(got, "-") {
		t.Fatalf("slug %q ends with a dash", got)
	}
}
