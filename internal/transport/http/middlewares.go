package http

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Si luego quieres, aquí puedes generar un ID y ponerlo en header
		// w.Header().Set("X-Request-Id", "...")
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

// methodOverrideMiddleware permite HTML forms con POST simular PUT/DELETE/PATCH usando:
// - header: X-HTTP-Method-Override
// - o query: ?_method=PUT
func methodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			override := r.Header.Get("X-HTTP-Method-Override")
			if override == "" {
				override = r.URL.Query().Get("_method")
			}
			override = strings.ToUpper(strings.TrimSpace(override))
			if override == http.MethodPut || override == http.MethodDelete || override == http.MethodPatch {
				r.Method = override
			}
		}
		next.ServeHTTP(w, r)
	})
}
