package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/jonradoff/flipbook/internal/database"
	"github.com/jonradoff/flipbook/internal/models"
	"github.com/jonradoff/flipbook/internal/storage"
)

type ViewerHandler struct {
	db      *database.DB
	storage *storage.Storage
	tmpl    *template.Template
	baseURL string
	uiLang  string // "auto", "ko" or "en"
}

func NewViewerHandler(db *database.DB, store *storage.Storage, tmpl *template.Template, baseURL, uiLang string) *ViewerHandler {
	return &ViewerHandler{db: db, storage: store, tmpl: tmpl, baseURL: baseURL, uiLang: uiLang}
}

func (h *ViewerHandler) View(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	fb, err := h.db.GetFlipbookBySlug(slug)
	if err != nil {
		http.Error(w, "Flipbook not found", 404)
		return
	}

	// Title-only guess until page text is available (refined below when ready)
	lang := detectLang(fb.Title, nil)

	switch fb.Status {
	case models.StatusReady:
		// continue below
	case models.StatusRegenerating, models.StatusConverting, models.StatusPending:
		h.tmpl.ExecuteTemplate(w, "viewer_wait", map[string]interface{}{
			"Flipbook": fb,
			"BaseURL":  h.baseURL,
			"Lang":     lang,
			"T":        uiFor(resolveUILang(h.uiLang, lang)),
		})
		return
	default:
		http.Error(w, "Flipbook not found", 404)
		return
	}

	go h.db.RecordView(fb.ID, r.Referer(), r.UserAgent())

	pageFmt := h.storage.DetectPageFormat(fb.ID)
	var pages []string
	var thumbs []string
	for i := 1; i <= fb.PageCount; i++ {
		pages = append(pages, h.storage.PageImageURL(fb.ID, pageFmt, i))
		thumbs = append(thumbs, h.storage.ThumbImageURL(fb.ID, pageFmt, i))
	}

	// Load extracted text for search and SEO (nil if unavailable)
	pageTexts := h.storage.LoadPageTexts(fb.ID)
	pageTextsJSON, _ := json.Marshal(pageTexts)
	lang = detectLang(fb.Title, pageTexts)

	// Build SEO description from first few pages of text
	metaDesc := fb.Description
	if metaDesc == "" {
		metaDesc = buildDescription(fb.Title, pageTexts)
	}

	// Canonical URL and OG image
	canonicalURL := h.baseURL + "/v/" + url.PathEscape(fb.Slug)
	var ogImage string
	if fb.PageCount > 0 {
		ogImage = h.baseURL + h.storage.PageImageURL(fb.ID, pageFmt, 1)
	}

	h.tmpl.ExecuteTemplate(w, "viewer", map[string]interface{}{
		"Flipbook":      fb,
		"Pages":         pages,
		"Thumbs":        thumbs,
		"PageTexts":     pageTexts,
		"PageTextsJSON": template.JS(pageTextsJSON),
		"BaseURL":       h.baseURL,
		"CanonicalURL":  canonicalURL,
		"OGImage":       ogImage,
		"MetaDesc":      metaDesc,
		"EmbedCode":     embedCode(h.baseURL, fb.Slug),
		"Lang":          lang,
		"T":             uiFor(resolveUILang(h.uiLang, lang)),
		// Hidden SEO headings are content, so they follow the content language
		"SlideLabel": uiFor(lang)["slide"],
	})
}

// maxDescriptionRunes caps the meta description length in characters (not
// bytes, so Korean text is not cut mid-character).
const maxDescriptionRunes = 300

// buildDescription generates a meta description from the flipbook title and page text.
func buildDescription(title string, pageTexts []string) string {
	var parts []string
	charCount := 0
	for _, text := range pageTexts {
		text = strings.Join(strings.Fields(text), " ")
		if text == "" {
			continue
		}
		if remaining := maxDescriptionRunes - charCount; utf8.RuneCountInString(text) > remaining {
			if remaining > 20 {
				cut, _ := truncateRunes(text, remaining)
				parts = append(parts, cut+"...")
			}
			break
		}
		parts = append(parts, text)
		charCount += utf8.RuneCountInString(text) + 1
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%s — interactive flipbook with full-page presentation slides.", title)
	}
	return strings.Join(parts, " ")
}
