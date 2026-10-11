package stats

import (
	"testing"
	"time"

	"github.com/jonradoff/flipbook/internal/database"
)

const (
	chromeWin = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0 Safari/537.36"
	iphone    = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1"
)

func TestIsBot(t *testing.T) {
	bots := []string{"", "Googlebot/2.1 (+http://www.google.com/bot.html)", "facebookexternalhit/1.1",
		"Mozilla/5.0 (compatible; Yeti/1.1; +https://naver.me/spd)", "kakaotalk-scrap/1.0", "curl/8.4.0", "python-requests/2.31"}
	for _, ua := range bots {
		if !IsBot(ua) {
			t.Errorf("IsBot(%q) = false, want true", ua)
		}
	}
	for _, ua := range []string{chromeWin, iphone} {
		if IsBot(ua) {
			t.Errorf("IsBot(%q) = true, want false", ua)
		}
	}
}

func TestIsMobile(t *testing.T) {
	if !IsMobile(iphone) || IsMobile(chromeWin) {
		t.Fatal("IsMobile misclassified iPhone/desktop")
	}
}

func TestReferrerHost(t *testing.T) {
	tests := []struct{ ref, want string }{
		{"", ""},
		{"not a url", ""},
		{"https://blog.naver.com/readmaster/123", "blog.naver.com"},
		{"https://www.readmaster.kr/", "readmaster.kr"},
		{"https://flipbook-revised.fly.dev/", "이 사이트"},
	}
	for _, tt := range tests {
		if got := ReferrerHost(tt.ref, "flipbook-revised.fly.dev", "이 사이트"); got != tt.want {
			t.Errorf("ReferrerHost(%q) = %q, want %q", tt.ref, got, tt.want)
		}
	}
}

func TestParsePeriod(t *testing.T) {
	now := time.Date(2026, 10, 11, 1, 30, 0, 0, time.UTC) // 10:30 KST, Oct 11
	day := func(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 0, 0, 0, 0, KST) }
	tests := []struct {
		key      string
		from, to time.Time
	}{
		{"30d", day(2026, 9, 12), day(2026, 10, 12)},
		{"this-month", day(2026, 10, 1), day(2026, 10, 12)},
		{"last-month", day(2026, 9, 1), day(2026, 10, 1)},
		{"bogus", day(2026, 9, 12), day(2026, 10, 12)},
	}
	for _, tt := range tests {
		p := ParsePeriod(tt.key, now)
		if !p.From.Equal(tt.from) || !p.To.Equal(tt.to) {
			t.Errorf("ParsePeriod(%q) = %v..%v, want %v..%v", tt.key, p.From, p.To, tt.from, tt.to)
		}
	}
}

func TestSummarize(t *testing.T) {
	p := Period{From: time.Date(2026, 10, 1, 0, 0, 0, 0, KST), To: time.Date(2026, 10, 4, 0, 0, 0, 0, KST)}
	at := func(d, h int) time.Time { return time.Date(2026, 10, d, h, 0, 0, 0, KST) }
	views := []database.ViewRecord{
		{FlipbookID: "a", ViewedAt: at(1, 9), Referrer: "https://blog.naver.com/x", UserAgent: iphone},
		{FlipbookID: "a", ViewedAt: at(1, 23), Referrer: "https://blog.naver.com/y", UserAgent: chromeWin},
		{FlipbookID: "b", ViewedAt: at(3, 8), Referrer: "", UserAgent: iphone},
		// 23:00 Sep 30 KST is before the period starts.
		{FlipbookID: "a", ViewedAt: time.Date(2026, 9, 30, 23, 0, 0, 0, KST), UserAgent: chromeWin},
		// 00:30 Oct 4 KST is already past the end.
		{FlipbookID: "b", ViewedAt: time.Date(2026, 10, 3, 15, 30, 0, 0, time.UTC), UserAgent: chromeWin},
	}
	s := Summarize(views, p, "flipbook-revised.fly.dev", "이 사이트", 5)

	if s.Total != 3 {
		t.Fatalf("Total = %d, want 3", s.Total)
	}
	if s.ByFlipbook["a"] != 2 || s.ByFlipbook["b"] != 1 {
		t.Fatalf("ByFlipbook = %v", s.ByFlipbook)
	}
	if len(s.Daily) != 3 || s.Daily[0].Count != 2 || s.Daily[1].Count != 0 || s.Daily[2].Count != 1 {
		t.Fatalf("Daily = %+v, want [2 0 1]", s.Daily)
	}
	if s.Referrers[0] != (Named{Name: "blog.naver.com", Count: 2}) || s.Referrers[1] != (Named{Name: "", Count: 1}) {
		t.Fatalf("Referrers = %+v", s.Referrers)
	}
	if s.Mobile != 2 || s.Desktop != 1 || s.MobilePercent() != 66 {
		t.Fatalf("Mobile=%d Desktop=%d pct=%d", s.Mobile, s.Desktop, s.MobilePercent())
	}
	if s.MaxDaily() != 2 {
		t.Fatalf("MaxDaily = %d, want 2", s.MaxDaily())
	}
}

func TestSummarizeEmpty(t *testing.T) {
	p := ParsePeriod("30d", time.Now())
	s := Summarize(nil, p, "", "", 5)
	if s.Total != 0 || len(s.Daily) != 30 || s.MobilePercent() != 0 || s.MaxDaily() != 1 {
		t.Fatalf("empty summary = total %d, %d days, pct %d, max %d", s.Total, len(s.Daily), s.MobilePercent(), s.MaxDaily())
	}
}
