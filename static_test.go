package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticCacheControl(t *testing.T) {
	year := "public, max-age=31536000, immutable"
	tests := []struct{ path, query, want string }{
		{"/static/css/style.css", "v=8", year},
		{"/static/js/viewer.js", "v=12", year},
		{"/static/js/viewer.js", "x=1&v=12", year},
		{"/static/js/viewer.js", "dev=1", "public, max-age=3600"},
		{"/static/fonts/sametone/SameToneText-400-latin.woff2", "", "public, max-age=2592000"},
		{"/static/img/landing/hero-flip.GIF", "", "public, max-age=604800"},
		{"/static/img/favicon.svg", "", "public, max-age=604800"},
		{"/static/fonts/sametone/sametone.css", "", "public, max-age=3600"},
		{"/static/fonts/sametone/OFL.txt", "", "public, max-age=3600"},
	}
	for _, tt := range tests {
		if got := staticCacheControl(tt.path, tt.query); got != tt.want {
			t.Errorf("staticCacheControl(%q, %q) = %q, want %q", tt.path, tt.query, got, tt.want)
		}
	}
}

func TestWithStaticCacheSetsHeader(t *testing.T) {
	h := withStaticCache(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) }))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/static/css/viewer.css?v=10", nil))
	if got := rec.Header().Get("Cache-Control"); got != "public, max-age=31536000, immutable" {
		t.Fatalf("Cache-Control = %q", got)
	}
}

func TestWithStaticCacheSkipsErrors(t *testing.T) {
	h := withStaticCache(http.NotFoundHandler())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/static/css/missing.css?v=10", nil))
	if rec.Code != http.StatusNotFound || rec.Header().Get("Cache-Control") != "" {
		t.Fatalf("404 response: code %d, Cache-Control %q; want no caching", rec.Code, rec.Header().Get("Cache-Control"))
	}
}
