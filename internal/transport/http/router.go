package http

import (
    "net/http"

    "github.com/gorilla/mux"
)

type Server struct {
    Router *mux.Router
    h      *Handlers
}

func NewServer(h *Handlers) *Server {
    r := mux.NewRouter()
    s := &Server{Router: r, h: h}

    r.Use(requestIDMiddleware)
    r.Use(loggingMiddleware)
    r.Use(methodOverrideMiddleware)

    // Health
    r.HandleFunc("/health", h.Health).Methods(http.MethodGet)

    // UI
    r.HandleFunc("/", h.UIHome).Methods(http.MethodGet)
    r.HandleFunc("/ui/users", h.UIUsers).Methods(http.MethodGet, http.MethodPost)
    r.HandleFunc("/ui/books", h.UIBooks).Methods(http.MethodGet, http.MethodPost)
    r.HandleFunc("/ui/books/search", h.UIBookSearch).Methods(http.MethodGet)
    r.HandleFunc("/ui/books/{id:[0-9]+}", h.UIBookDetail).Methods(http.MethodGet)
    r.HandleFunc("/ui/access", h.UIAccess).Methods(http.MethodPost)

    // API Users
    r.HandleFunc("/api/users", h.APIListUsers).Methods(http.MethodGet)
    r.HandleFunc("/api/users", h.APICreateUser).Methods(http.MethodPost)
    r.HandleFunc("/api/users/{id:[0-9]+}", h.APIGetUser).Methods(http.MethodGet)
    r.HandleFunc("/api/users/{id:[0-9]+}", h.APIUpdateUser).Methods(http.MethodPut)
    r.HandleFunc("/api/users/{id:[0-9]+}", h.APIDeleteUser).Methods(http.MethodDelete)

    // API Books
    r.HandleFunc("/api/books", h.APIListBooks).Methods(http.MethodGet)
    r.HandleFunc("/api/books", h.APICreateBook).Methods(http.MethodPost)
    r.HandleFunc("/api/books/search", h.APISearchBooks).Methods(http.MethodGet)
    r.HandleFunc("/api/books/{id:[0-9]+}", h.APIGetBook).Methods(http.MethodGet)
    r.HandleFunc("/api/books/{id:[0-9]+}", h.APIUpdateBook).Methods(http.MethodPut)
    r.HandleFunc("/api/books/{id:[0-9]+}", h.APIDeleteBook).Methods(http.MethodDelete)
    r.HandleFunc("/api/books/{id:[0-9]+}/stats", h.APIBookStats).Methods(http.MethodGet)

    // API Access
    r.HandleFunc("/api/access", h.APIRecordAccess).Methods(http.MethodPost)

    return s
}
