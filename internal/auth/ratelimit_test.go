package auth

import (
	"html/template"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoginLimiter(t *testing.T) {
	now := time.Date(2026, 10, 11, 9, 0, 0, 0, time.UTC)
	l := NewLimiter(3, 10*time.Minute)
	l.now = func() time.Time { return now }

	for i := 0; i < 3; i++ {
		if blocked, _ := l.Blocked("1.1.1.1"); blocked {
			t.Fatalf("blocked after %d failures, want 3 allowed", i)
		}
		l.Fail("1.1.1.1")
		now = now.Add(time.Minute)
	}
	blocked, wait := l.Blocked("1.1.1.1")
	if !blocked || wait != 7*time.Minute {
		t.Fatalf("after 3 failures: blocked=%v wait=%v, want true 7m", blocked, wait)
	}
	if b, _ := l.Blocked("2.2.2.2"); b {
		t.Fatal("other IPs must not be blocked")
	}

	now = now.Add(7 * time.Minute) // oldest failure leaves the window
	if b, _ := l.Blocked("1.1.1.1"); b {
		t.Fatal("still blocked after the window passed")
	}

	l.Fail("1.1.1.1")
	l.Reset("1.1.1.1") // successful login clears the count
	if len(l.failures) != 0 {
		t.Fatalf("Reset left %v", l.failures)
	}
}

func TestClientIP(t *testing.T) {
	r := httptest.NewRequest("POST", "/login", nil)
	r.RemoteAddr = "10.0.0.5:4321"
	r.Header.Set("Fly-Client-IP", "203.0.113.9")

	onFly = false
	if got := ClientIP(r); got != "10.0.0.5" {
		t.Fatalf("off Fly: %q, want RemoteAddr host (header must be ignored)", got)
	}
	onFly = true
	defer func() { onFly = false }()
	if got := ClientIP(r); got != "203.0.113.9" {
		t.Fatalf("on Fly: %q, want Fly-Client-IP", got)
	}
}

func TestLoginSubmitBlocked(t *testing.T) {
	tmpl := template.Must(template.New("login").Parse(`{{define "login"}}{{.Error}}{{end}}`))
	a := &Auth{tmpl: tmpl, uiLang: "en", limit: NewLimiter(1, time.Minute)}
	a.limit.Fail("10.0.0.5")

	r := httptest.NewRequest("POST", "/login", strings.NewReader("password=x"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.RemoteAddr = "10.0.0.5:1"
	rec := httptest.NewRecorder()
	a.LoginSubmit(rec, r) // db is nil: a blocked request must not reach the password check

	if rec.Code != 429 || rec.Header().Get("Retry-After") == "" || !strings.Contains(rec.Body.String(), "Too many") {
		t.Fatalf("code=%d retry=%q body=%q", rec.Code, rec.Header().Get("Retry-After"), rec.Body.String())
	}
}
