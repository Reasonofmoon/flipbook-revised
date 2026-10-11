package handlers

import (
	"crypto/subtle"
	"net/http"
	"net/url"

	"github.com/jonradoff/flipbook/internal/models"
)

// accessQuery is the "?k=..." suffix a link-only flipbook needs, or "".
func accessQuery(fb *models.Flipbook) string {
	if fb.IsPublic || fb.AccessKey == "" {
		return ""
	}
	return "?k=" + url.QueryEscape(fb.AccessKey)
}

// viewerPath and embedPath are the only place viewer/embed links are built,
// so a link-only flipbook never loses its key.
func viewerPath(fb *models.Flipbook) string {
	return "/v/" + url.PathEscape(fb.Slug) + accessQuery(fb)
}

func embedPath(fb *models.Flipbook) string {
	return "/embed/" + url.PathEscape(fb.Slug) + accessQuery(fb)
}

// canView reports whether the request may see the flipbook: always when it
// is public, otherwise only with the matching access key.
func canView(fb *models.Flipbook, r *http.Request) bool {
	if fb.IsPublic {
		return true
	}
	given := r.URL.Query().Get("k")
	return fb.AccessKey != "" && given != "" &&
		subtle.ConstantTimeCompare([]byte(given), []byte(fb.AccessKey)) == 1
}

// pageAccessKey is the key the viewer page passes to its share dialog: the
// flipbook's key when it is link-only, "" when it is public.
func pageAccessKey(fb *models.Flipbook) string {
	if fb.IsPublic {
		return ""
	}
	return fb.AccessKey
}
