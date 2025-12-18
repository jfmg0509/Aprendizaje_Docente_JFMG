package http

import (
    "context"
    "log"
    "math/rand"
    "net/http"
    "time"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rid := randomID()
        ctx := context.WithValue(r.Context(), requestIDKey, rid)
        w.Header().Set("X-Request-Id", rid)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        rid, _ := r.Context().Value(requestIDKey).(string)
        log.Printf("[%s] %s %s (%s)", rid, r.Method, r.URL.Path, time.Since(start))
    })
}

// Permite enviar PUT/DELETE desde forms HTML: <input name="_method" value="PUT">
func methodOverrideMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            if err := r.ParseForm(); err == nil {
                if m := r.FormValue("_method"); m != "" {
                    r.Method = m
                }
            }
        }
        next.ServeHTTP(w, r)
    })
}

func randomID() string {
    rand.Seed(time.Now().UnixNano())
    const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, 10)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
