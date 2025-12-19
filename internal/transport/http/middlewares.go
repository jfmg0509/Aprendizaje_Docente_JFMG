package http

import (
	"log"
	"net/http"
	"time"
)

// Logging básico de requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

// Middleware para sobrescribir método vía formulario (_method)
func methodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if m := r.FormValue("_method"); m != "" {
				r.Method = m
			}
		}
		next.ServeHTTP(w, r)
	})
}
