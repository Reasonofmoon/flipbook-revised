package handlers

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

// LandingHandler serves the public funnel page at /.
type LandingHandler struct {
	tmpl     *template.Template
	kakaoURL string
	demoSlug string
}

func NewLandingHandler(tmpl *template.Template, kakaoURL, demoSlug string) *LandingHandler {
	return &LandingHandler{tmpl: tmpl, kakaoURL: safeExternalURL(kakaoURL), demoSlug: demoSlug}
}

func (h *LandingHandler) Index(w http.ResponseWriter, r *http.Request) {
	demoURL := ""
	if h.demoSlug != "" {
		demoURL = "/v/" + url.PathEscape(h.demoSlug)
	}
	h.tmpl.ExecuteTemplate(w, "landing", map[string]interface{}{
		"KakaoURL": h.kakaoURL,
		"DemoURL":  demoURL,
	})
}

// safeExternalURL keeps only absolute http(s) links; anything else disables the CTA.
func safeExternalURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		return ""
	}
	return u.String()
}
