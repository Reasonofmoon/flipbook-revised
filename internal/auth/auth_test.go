package auth

import (
	"net/http/httptest"
	"testing"
)

func TestSecureCookies(t *testing.T) {
	for url, want := range map[string]bool{
		"https://flipbook-revised.fly.dev": true,
		"HTTPS://Example.com":              true,
		" https://x.test ":                 true,
		"http://localhost:8080":            false,
		"":                                 false,
	} {
		if got := SecureCookies(url); got != want {
			t.Errorf("SecureCookies(%q) = %v, want %v", url, got, want)
		}
	}
}

func TestLogoutCookieFlags(t *testing.T) {
	for _, secure := range []bool{true, false} {
		a := &Auth{secure: secure} // no session cookie on the request, so db is not touched
		rec := httptest.NewRecorder()
		a.Logout(rec, httptest.NewRequest("POST", "/logout", nil))
		c := rec.Result().Cookies()
		if len(c) != 1 || c[0].Name != SessionCookieName || c[0].MaxAge != -1 || !c[0].HttpOnly || c[0].Secure != secure {
			t.Fatalf("secure=%v: logout cookie = %+v", secure, c)
		}
	}
}
