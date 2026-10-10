package handlers

import "testing"

// Every language must define every key that English defines, or the
// template renders an empty string for the missing text.
func TestUITextComplete(t *testing.T) {
	for lang, strs := range uiText {
		for key := range uiText["en"] {
			if _, ok := strs[key]; !ok {
				t.Errorf("uiText[%q] is missing key %q", lang, key)
			}
		}
		for key := range strs {
			if _, ok := uiText["en"][key]; !ok {
				t.Errorf("uiText[%q] has key %q not defined in en", lang, key)
			}
		}
	}
}

func TestUIForFallsBackToEnglish(t *testing.T) {
	if got := uiFor("fr")["copy"]; got != "Copy" {
		t.Fatalf(`uiFor("fr")["copy"] = %q, want "Copy"`, got)
	}
	if got := uiFor("ko")["copy"]; got != "복사" {
		t.Fatalf(`uiFor("ko")["copy"] = %q, want "복사"`, got)
	}
}

func TestResolveUILang(t *testing.T) {
	tests := []struct{ setting, content, want string }{
		{"auto", "en", "en"},
		{"auto", "ko", "ko"},
		{"ko", "en", "ko"}, // English-only deck shown with Korean UI
		{"en", "ko", "en"},
	}
	for _, tt := range tests {
		if got := resolveUILang(tt.setting, tt.content); got != tt.want {
			t.Errorf("resolveUILang(%q, %q) = %q, want %q", tt.setting, tt.content, got, tt.want)
		}
	}
}
