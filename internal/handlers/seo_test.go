package handlers

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestDetectLang(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		pageTexts []string
		want      string
	}{
		{"english deck", "Quarterly Review", []string{"Revenue grew 12% this quarter."}, "en"},
		{"korean deck", "수능 영어 독해", []string{"핵심 어휘 정리"}, "ko"},
		{"english lesson with korean glosses", "Vocabulary", []string{"abundant 풍부한 / scarce 부족한 / vivid 생생한"}, "ko"},
		{"one korean word in english text", "Seoul Trip", []string{"We visited Gyeongbokgung palace and ate 김밥 for lunch with friends."}, "en"},
		{"korean title, no text yet", "모의고사 해설", nil, "ko"},
		{"no letters", "2024", []string{"123 456"}, "en"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectLang(tt.title, tt.pageTexts); got != tt.want {
				t.Errorf("detectLang(%q) = %q, want %q", tt.title, got, tt.want)
			}
		})
	}
}

func TestBuildDescriptionKeepsValidUTF8(t *testing.T) {
	long := strings.Repeat("한국어 교재 설명문 ", 60) // well over the limit
	desc := buildDescription("제목", []string{long})

	if !utf8.ValidString(desc) {
		t.Fatal("description contains invalid UTF-8 (cut mid-character)")
	}
	if !strings.HasSuffix(desc, "...") {
		t.Fatalf("expected truncated description to end with ..., got %q", desc)
	}
	if n := utf8.RuneCountInString(desc); n > maxDescriptionRunes+3 {
		t.Fatalf("description has %d runes, want <= %d", n, maxDescriptionRunes+3)
	}
}
