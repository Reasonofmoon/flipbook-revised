package auth

import (
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	loginMaxFailures = 5
	loginWindow      = 15 * time.Minute
)

// loginLimiter tracks failed logins per client IP in memory. The app runs on a
// single machine, so process memory is enough; a restart simply resets it.
type loginLimiter struct {
	mu       sync.Mutex
	failures map[string][]time.Time
	max      int
	window   time.Duration
	now      func() time.Time
}

func newLoginLimiter(max int, window time.Duration) *loginLimiter {
	return &loginLimiter{failures: map[string][]time.Time{}, max: max, window: window, now: time.Now}
}

// recent drops failures older than the window. Caller holds mu.
func (l *loginLimiter) recent(ip string) []time.Time {
	cutoff := l.now().Add(-l.window)
	kept := l.failures[ip][:0]
	for _, t := range l.failures[ip] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) == 0 {
		delete(l.failures, ip)
		return nil
	}
	l.failures[ip] = kept
	return kept
}

// Blocked reports whether ip has used up its attempts, and how long until the
// oldest failure leaves the window.
func (l *loginLimiter) Blocked(ip string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f := l.recent(ip)
	if len(f) < l.max {
		return false, 0
	}
	return true, f[0].Add(l.window).Sub(l.now())
}

func (l *loginLimiter) Fail(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.failures[ip] = append(l.recent(ip), l.now())
	// Bound memory if many IPs hammer the endpoint: prune expired entries.
	if len(l.failures) > 10000 {
		for k := range l.failures {
			l.recent(k)
		}
	}
}

func (l *loginLimiter) Reset(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.failures, ip)
}

// onFly is true when running on Fly.io, where the proxy sets Fly-Client-IP.
var onFly = os.Getenv("FLY_APP_NAME") != ""

// clientIP returns the caller's IP. Behind Fly's proxy RemoteAddr is the
// proxy, so the Fly-Client-IP header is used there; elsewhere the header is
// ignored because any client could send it.
func clientIP(r *http.Request) string {
	if onFly {
		if ip := r.Header.Get("Fly-Client-IP"); ip != "" {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
