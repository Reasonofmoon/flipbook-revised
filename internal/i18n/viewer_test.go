package i18n

import "testing"

// Every language must define every key that English defines (and no extra
// keys), or a template renders an empty string for the missing text.
func TestTablesComplete(t *testing.T) {
	tables := map[string]map[string]Strings{"viewer": viewerText, "admin": adminText}
	for name, table := range tables {
		for lang, strs := range table {
			for key := range table["en"] {
				if _, ok := strs[key]; !ok {
					t.Errorf("%s[%q] is missing key %q", name, lang, key)
				}
			}
			for key := range strs {
				if _, ok := table["en"][key]; !ok {
					t.Errorf("%s[%q] has key %q not defined in en", name, lang, key)
				}
			}
		}
	}
}

func TestLookupFallsBackToEnglish(t *testing.T) {
	if got := Viewer("fr")["copy"]; got != "Copy" {
		t.Fatalf(`Viewer("fr")["copy"] = %q, want "Copy"`, got)
	}
	if got := Viewer("ko")["copy"]; got != "복사" {
		t.Fatalf(`Viewer("ko")["copy"] = %q, want "복사"`, got)
	}
	if got := Admin("fr")["logout"]; got != "Logout" {
		t.Fatalf(`Admin("fr")["logout"] = %q, want "Logout"`, got)
	}
}

func TestResolveViewer(t *testing.T) {
	tests := []struct{ setting, content, want string }{
		{"auto", "en", "en"},
		{"auto", "ko", "ko"},
		{"ko", "en", "ko"}, // English-only deck shown with Korean UI
		{"en", "ko", "en"},
	}
	for _, tt := range tests {
		if got := ResolveViewer(tt.setting, tt.content); got != tt.want {
			t.Errorf("ResolveViewer(%q, %q) = %q, want %q", tt.setting, tt.content, got, tt.want)
		}
	}
}

func TestResolveAdmin(t *testing.T) {
	tests := []struct{ setting, accept, want string }{
		{"ko", "en-US,en;q=0.9", "ko"},
		{"en", "ko-KR,ko;q=0.9", "en"},
		{"auto", "ko-KR,ko;q=0.9,en-US;q=0.8", "ko"},
		{"auto", "en-US,en;q=0.9,ko;q=0.8", "en"},
		{"auto", "ja-JP,ko;q=0.8", "ko"}, // first supported language wins
		{"auto", "fr-FR", "en"},
		{"auto", "", "en"},
	}
	for _, tt := range tests {
		if got := ResolveAdmin(tt.setting, tt.accept); got != tt.want {
			t.Errorf("ResolveAdmin(%q, %q) = %q, want %q", tt.setting, tt.accept, got, tt.want)
		}
	}
}
