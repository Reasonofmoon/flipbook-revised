package handlers

import "testing"

func TestSafeExternalURL(t *testing.T) {
	tests := map[string]string{
		"https://pf.kakao.com/_abcd/chat": "https://pf.kakao.com/_abcd/chat",
		"  https://pf.kakao.com/_abcd  ":  "https://pf.kakao.com/_abcd",
		"http://pf.kakao.com/_abcd":       "http://pf.kakao.com/_abcd",
		"":                                "",
		"pf.kakao.com/_abcd":              "", // no scheme
		"javascript:alert(1)":             "",
		"data:text/html,<script>":         "",
		"/relative/path":                  "",
	}
	for in, want := range tests {
		if got := safeExternalURL(in); got != want {
			t.Errorf("safeExternalURL(%q) = %q, want %q", in, got, want)
		}
	}
}
