package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jonradoff/flipbook/internal/models"
)

func TestCanView(t *testing.T) {
	public := &models.Flipbook{Slug: "unit-3", IsPublic: true, AccessKey: "abc123"}
	private := &models.Flipbook{Slug: "unit-3", IsPublic: false, AccessKey: "s3cret-KEY_x"}
	keyless := &models.Flipbook{Slug: "unit-3", IsPublic: false} // misconfigured: never viewable

	tests := []struct {
		name string
		fb   *models.Flipbook
		url  string
		want bool
	}{
		{"public, no key", public, "/v/unit-3", true},
		{"public, any key", public, "/v/unit-3?k=wrong", true},
		{"private, right key", private, "/v/unit-3?k=s3cret-KEY_x", true},
		{"private, right key with page", private, "/v/unit-3?page=3&k=s3cret-KEY_x", true},
		{"private, no key", private, "/v/unit-3", false},
		{"private, wrong key", private, "/v/unit-3?k=s3cret-KEY_y", false},
		{"private, key prefix", private, "/v/unit-3?k=s3cret", false},
		{"private, empty key", private, "/v/unit-3?k=", false},
		{"private without stored key", keyless, "/v/unit-3?k=", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", tt.url, nil)
			if got := canView(tt.fb, r); got != tt.want {
				t.Errorf("canView = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinksCarryKeyOnlyWhenPrivate(t *testing.T) {
	public := &models.Flipbook{Slug: "수능-영어", IsPublic: true, AccessKey: "abc"}
	private := &models.Flipbook{Slug: "수능-영어", IsPublic: false, AccessKey: "a+b/c"}

	if got := viewerPath(public); strings.Contains(got, "?k=") {
		t.Errorf("public viewerPath %q leaks a key", got)
	}
	if got, want := viewerPath(private), "/v/%EC%88%98%EB%8A%A5-%EC%98%81%EC%96%B4?k=a%2Bb%2Fc"; got != want {
		t.Errorf("private viewerPath = %q, want %q", got, want)
	}
	if got := embedPath(private); !strings.HasPrefix(got, "/embed/") || !strings.HasSuffix(got, "?k=a%2Bb%2Fc") {
		t.Errorf("private embedPath = %q", got)
	}
	if code := embedCode("https://x.test", private); !strings.Contains(code, `src="https://x.test/embed/`) || !strings.Contains(code, "?k=a%2Bb%2Fc") {
		t.Errorf("embedCode misses the key: %s", code)
	}
	if pageAccessKey(public) != "" || pageAccessKey(private) != "a+b/c" {
		t.Error("pageAccessKey should expose the key only for private flipbooks")
	}
}
