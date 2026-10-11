package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jonradoff/flipbook/internal/auth"
)

func TestAPIAuthLockout(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	h := apiAuth("secret-key", auth.NewLimiter(2, time.Minute))(ok)

	call := func(ip, token string) *httptest.ResponseRecorder {
		r := httptest.NewRequest("GET", "/api/flipbooks", nil)
		r.RemoteAddr = ip + ":1234"
		if token != "" {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec
	}

	if c := call("10.0.0.1", "secret-key").Code; c != 200 {
		t.Fatalf("right key: %d, want 200", c)
	}
	for i := 0; i < 2; i++ {
		if c := call("10.0.0.2", "wrong").Code; c != 401 {
			t.Fatalf("wrong key #%d: %d, want 401", i+1, c)
		}
	}
	rec := call("10.0.0.2", "secret-key") // locked: even the right key is refused
	if rec.Code != 429 || rec.Header().Get("Retry-After") == "" {
		t.Fatalf("locked IP: %d retry=%q, want 429 with Retry-After", rec.Code, rec.Header().Get("Retry-After"))
	}
	if c := call("10.0.0.1", "secret-key").Code; c != 200 {
		t.Fatalf("other IP affected by lockout: %d", c)
	}
	if c := call("10.0.0.3", "").Code; c != 401 {
		t.Fatalf("missing header: %d, want 401", c)
	}
}
