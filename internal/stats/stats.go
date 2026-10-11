// Package stats turns raw flipbook views into report numbers: totals, a daily
// series, top referrers and the phone/desktop split. All functions are pure.
package stats

import (
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jonradoff/flipbook/internal/database"
)

// KST is Korea Standard Time. Korea has no DST, and the runtime image has no
// tzdata, so a fixed zone is used instead of time.LoadLocation.
var KST = time.FixedZone("KST", 9*60*60)

var botRe = regexp.MustCompile(`(?i)bot|crawl|spider|slurp|preview|facebookexternalhit|kakaotalk-scrap|daum|yeti|headless|curl|wget|python-requests|go-http-client`)

// IsBot reports whether a user agent belongs to a crawler, link-preview
// fetcher or script rather than a person. An empty UA counts as a bot.
func IsBot(userAgent string) bool {
	return strings.TrimSpace(userAgent) == "" || botRe.MatchString(userAgent)
}

var mobileRe = regexp.MustCompile(`(?i)mobile|android|iphone|ipad|ipod`)

// IsMobile reports whether a user agent is a phone or tablet.
func IsMobile(userAgent string) bool { return mobileRe.MatchString(userAgent) }

// ReferrerHost returns the referring site's host without "www.", "" for a
// direct visit, or selfLabel when the visit came from this site itself.
func ReferrerHost(referrer, selfHost, selfLabel string) string {
	u, err := url.Parse(strings.TrimSpace(referrer))
	if err != nil || u.Host == "" {
		return ""
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	if selfHost != "" && host == strings.TrimPrefix(strings.ToLower(selfHost), "www.") {
		return selfLabel
	}
	return host
}

// Period is a reporting window [From, To) in KST.
type Period struct {
	Key  string // "30d", "this-month" or "last-month"
	From time.Time
	To   time.Time
}

// ParsePeriod resolves a period key relative to now; unknown keys mean "30d".
func ParsePeriod(key string, now time.Time) Period {
	now = now.In(KST)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, KST)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, KST)
	switch key {
	case "this-month":
		return Period{Key: key, From: monthStart, To: today.AddDate(0, 0, 1)}
	case "last-month":
		return Period{Key: key, From: monthStart.AddDate(0, -1, 0), To: monthStart}
	default:
		return Period{Key: "30d", From: today.AddDate(0, 0, -29), To: today.AddDate(0, 0, 1)}
	}
}

// Day is the view count for one KST calendar day.
type Day struct {
	Date  time.Time
	Count int
}

// Named is a label with a count, e.g. a referrer host.
type Named struct {
	Name  string
	Count int
}

// Summary aggregates views over a period.
type Summary struct {
	Total      int
	Daily      []Day          // every day in the period, zero-filled
	ByFlipbook map[string]int // flipbook ID -> views
	Referrers  []Named        // most first; "" means direct visit
	Mobile     int
	Desktop    int
}

// MobilePercent is the share of phone/tablet views, rounded down.
func (s Summary) MobilePercent() int {
	if s.Total == 0 {
		return 0
	}
	return s.Mobile * 100 / s.Total
}

// MaxDaily is the largest daily count (at least 1, for chart scaling).
func (s Summary) MaxDaily() int {
	m := 1
	for _, d := range s.Daily {
		if d.Count > m {
			m = d.Count
		}
	}
	return m
}

// Summarize aggregates views that fall inside p. Referrers keeps the top n
// entries (n <= 0 keeps all). selfHost/selfLabel group visits from this site.
func Summarize(views []database.ViewRecord, p Period, selfHost, selfLabel string, n int) Summary {
	s := Summary{ByFlipbook: map[string]int{}}
	days := map[string]int{}
	refs := map[string]int{}
	for _, v := range views {
		t := v.ViewedAt.In(KST)
		if t.Before(p.From) || !t.Before(p.To) {
			continue
		}
		s.Total++
		s.ByFlipbook[v.FlipbookID]++
		days[t.Format("2006-01-02")]++
		refs[ReferrerHost(v.Referrer, selfHost, selfLabel)]++
		if IsMobile(v.UserAgent) {
			s.Mobile++
		} else {
			s.Desktop++
		}
	}
	for d := p.From; d.Before(p.To); d = d.AddDate(0, 0, 1) {
		s.Daily = append(s.Daily, Day{Date: d, Count: days[d.Format("2006-01-02")]})
	}
	for name, c := range refs {
		s.Referrers = append(s.Referrers, Named{Name: name, Count: c})
	}
	sort.Slice(s.Referrers, func(i, j int) bool {
		if s.Referrers[i].Count != s.Referrers[j].Count {
			return s.Referrers[i].Count > s.Referrers[j].Count
		}
		return s.Referrers[i].Name < s.Referrers[j].Name
	})
	if n > 0 && len(s.Referrers) > n {
		s.Referrers = s.Referrers[:n]
	}
	return s
}
