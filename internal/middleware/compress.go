package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	Writer *gzip.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to create gzip reader", http.StatusBadRequest)
				return
			}
			defer gz.Close()
			r.Body = io.NopCloser(gz)
		}

		acceptGzip := false
		if accept := r.Header.Get("Accept-Encoding"); accept != "" {
			acceptGzip = contains(accept, "gzip")
		}

		if acceptGzip {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")
			gz := gzip.NewWriter(w)
			defer gz.Close()

			wrapped := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
			next.ServeHTTP(wrapped, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func contains(s, substr string) bool {
	for _, part := range splitTrim(s) {
		if part == substr {
			return true
		}
	}
	return false
}

func splitTrim(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
