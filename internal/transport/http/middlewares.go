package http

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware simple de RequestID (mínimo)
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// No generamos UUID para no meter libs. Solo timestamp nano.
		rid := time.Now().UnixNano()
		w.Header().Set("X-Request-ID", fmtInt64(rid))
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

// Permite “override” con ?_method=PUT (útil si luego haces forms)
func methodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			m := r.URL.Query().Get("_method")
			if m != "" {
				r.Method = m
			}
		}
		next.ServeHTTP(w, r)
	})
}

func fmtInt64(v int64) string {
	// sin strconv.FormatInt para mantenerlo simple, pero ok usar strconv:
	// return strconv.FormatInt(v, 10)
	// Aquí lo hacemos directo:
	return fmt.Sprintf("%d", v)
}
