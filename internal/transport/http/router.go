package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(h *Handler) http.Handler {
	r := mux.NewRouter()

	// Middlewares (si existen en tu proyecto)
	// Si alguno NO existe, bórralo y listo.
	r.Use(requestIDMiddleware)
	r.Use(loggingMiddleware)
	r.Use(methodOverrideMiddleware)

	// ==============================
	// UI (HTML)
	// ==============================
	r.HandleFunc("/", h.uiHome).Methods(http.MethodGet)

	r.HandleFunc("/ui/users", h.uiUsersGET).Methods(http.MethodGet)
	r.HandleFunc("/ui/users", h.uiUsersPOST).Methods(http.MethodPost)

	r.HandleFunc("/ui/books", h.uiBooksGET).Methods(http.MethodGet)
	r.HandleFunc("/ui/books", h.uiBooksPOST).Methods(http.MethodPost)

	r.HandleFunc("/ui/books/search", h.uiBookSearchGET).Methods(http.MethodGet)
	r.HandleFunc("/ui/books/{id:[0-9]+}", h.uiBookDetailGET).Methods(http.MethodGet)

	r.HandleFunc("/ui/access", h.uiAccessPOST).Methods(http.MethodPost)

	// ==============================
	// API (JSON)
	// ==============================
	api := r.PathPrefix("/api").Subrouter()

	// USERS API
	api.HandleFunc("/users", h.apiCreateUser).Methods(http.MethodPost)
	api.HandleFunc("/users", h.apiListUsers).Methods(http.MethodGet)
	api.HandleFunc("/users/{id:[0-9]+}", h.apiGetUser).Methods(http.MethodGet)
	api.HandleFunc("/users/{id:[0-9]+}", h.apiUpdateUser).Methods(http.MethodPut)
	api.HandleFunc("/users/{id:[0-9]+}", h.apiDeleteUser).Methods(http.MethodDelete)

	// BOOKS API
	api.HandleFunc("/books", h.apiCreateBook).Methods(http.MethodPost)
	api.HandleFunc("/books", h.apiListBooks).Methods(http.MethodGet)
	api.HandleFunc("/books/search", h.apiSearchBooks).Methods(http.MethodGet)
	api.HandleFunc("/books/{id:[0-9]+}", h.apiGetBook).Methods(http.MethodGet)
	api.HandleFunc("/books/{id:[0-9]+}", h.apiUpdateBook).Methods(http.MethodPatch)
	api.HandleFunc("/books/{id:[0-9]+}", h.apiDeleteBook).Methods(http.MethodDelete)

	// ACCESS API
	api.HandleFunc("/access", h.apiRecordAccess).Methods(http.MethodPost)
	api.HandleFunc("/books/{id:[0-9]+}/stats", h.apiStatsByBook).Methods(http.MethodGet)

	// HEALTH
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	return r
}
