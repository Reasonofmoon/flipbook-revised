package main

import (
	"net/http"
	"path"
	"strings"
)

// staticCacheControl picks a Cache-Control value for a file under /static.
// Files referenced with a ?v= version change URL whenever they change, so
// they can be cached for a year; unversioned files get shorter lifetimes.
func staticCacheControl(urlPath, rawQuery string) string {
	if strings.Contains("&"+rawQuery, "&v=") {
		return "public, max-age=31536000, immutable"
	}
	switch strings.ToLower(path.Ext(urlPath)) {
	case ".woff2", ".woff", ".ttf":
		return "public, max-age=2592000" // 30 days; fonts are rebuilt rarely
	case ".gif", ".png", ".jpg", ".jpeg", ".webp", ".svg", ".ico":
		return "public, max-age=604800" // 7 days; landing GIFs may be re-recorded
	default:
		return "public, max-age=3600" // unversioned CSS/JS and text files
	}
}

// withStaticCache sets Cache-Control on successful static file responses.
// Errors (e.g. 404 for a mistyped path) are left uncached so a bad URL is
// not remembered for a year.
func withStaticCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&cacheWriter{ResponseWriter: w, value: staticCacheControl(r.URL.Path, r.URL.RawQuery)}, r)
	})
}

type cacheWriter struct {
	http.ResponseWriter
	value       string
	wroteHeader bool
}

func (c *cacheWriter) WriteHeader(status int) {
	if !c.wroteHeader {
		c.wroteHeader = true
		if status < 400 {
			c.Header().Set("Cache-Control", c.value)
		}
	}
	c.ResponseWriter.WriteHeader(status)
}

func (c *cacheWriter) Write(b []byte) (int, error) {
	if !c.wroteHeader {
		c.WriteHeader(http.StatusOK)
	}
	return c.ResponseWriter.Write(b)
}
