package http

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// requestIDMiddleware: agrega un request id simple para rastrear logs.
// (Sin dependencias externas)
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Si ya viene uno, lo respetamos.
		rid := r.Header.Get("X-Request-Id")
		if rid == "" {
			rid = time.Now().Format("20060102-150405.000000000")
		}

		// Guardamos en headers de respuesta para debug.
		w.Header().Set("X-Request-Id", rid)

		// Continuar
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware: log básico por request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Ejecuta siguiente handler
		next.ServeHTTP(w, r)
		// Log al final
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// methodOverrideMiddleware: permite simular PUT/DELETE desde formularios HTML.
// Soporta:
// - Header: X-HTTP-Method-Override
// - Query:  ?_method=PUT
// - Form:   _method=PUT   (solo si POST)
func methodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1) header
		if m := strings.TrimSpace(r.Header.Get("X-HTTP-Method-Override")); m != "" {
			r.Method = strings.ToUpper(m)
			next.ServeHTTP(w, r)
			return
		}

		// 2) query
		if m := strings.TrimSpace(r.URL.Query().Get("_method")); m != "" {
			r.Method = strings.ToUpper(m)
			next.ServeHTTP(w, r)
			return
		}

		// 3) form (solo tiene sentido en POST)
		if r.Method == http.MethodPost {
			_ = r.ParseForm()
			if m := strings.TrimSpace(r.FormValue("_method")); m != "" {
				r.Method = strings.ToUpper(m)
			}
		}

		next.ServeHTTP(w, r)
	})
}
