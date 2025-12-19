package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter arma el enrutador principal (UI + API) y devuelve el router listo.
func NewRouter(ui *UIHandler, api *APIHandler) *mux.Router {
	r := mux.NewRouter()

	// =========================
	// UI ROUTES (HTML)
	// =========================

	//r.HandleFunc("/ui/access", ui.AccessCreate).Methods(http.MethodPost)

	r.HandleFunc("/", ui.Home).Methods(http.MethodGet)

	r.HandleFunc("/ui/users", ui.UsersPage).Methods(http.MethodGet)
	r.HandleFunc("/ui/users", ui.UsersCreate).Methods(http.MethodPost)

	r.HandleFunc("/ui/books", ui.BooksPage).Methods(http.MethodGet)
	r.HandleFunc("/ui/books", ui.BooksCreate).Methods(http.MethodPost)

	r.HandleFunc("/ui/books/search", ui.BookSearch).Methods(http.MethodGet)
	r.HandleFunc("/ui/books/{id:[0-9]+}", ui.BookDetail).Methods(http.MethodGet)

	// Si ya tienes el handler de registrar accesos desde UI (POST /ui/access)
	// (si todavía no lo tienes, comenta esta línea por ahora)
	r.HandleFunc("/ui/access", ui.AccessCreate).Methods(http.MethodPost)

	// =========================
	// API ROUTES (JSON /api/*)
	// =========================
	// Si todavía no tienes APIHandler, puedes pasar nil en main y comentar este bloque.
	if api != nil {
		apiRouter := r.PathPrefix("/api").Subrouter()

		// Users API
		apiRouter.HandleFunc("/users", api.UsersList).Methods(http.MethodGet)
		apiRouter.HandleFunc("/users", api.UsersCreate).Methods(http.MethodPost)
		apiRouter.HandleFunc("/users/{id:[0-9]+}", api.UsersGet).Methods(http.MethodGet)
		apiRouter.HandleFunc("/users/{id:[0-9]+}", api.UsersUpdate).Methods(http.MethodPut)
		apiRouter.HandleFunc("/users/{id:[0-9]+}", api.UsersDelete).Methods(http.MethodDelete)

		// Books API
		apiRouter.HandleFunc("/books", api.BooksList).Methods(http.MethodGet)
		apiRouter.HandleFunc("/books", api.BooksCreate).Methods(http.MethodPost)
		apiRouter.HandleFunc("/books/{id:[0-9]+}", api.BooksGet).Methods(http.MethodGet)
		apiRouter.HandleFunc("/books/{id:[0-9]+}", api.BooksUpdate).Methods(http.MethodPut)
		apiRouter.HandleFunc("/books/{id:[0-9]+}", api.BooksDelete).Methods(http.MethodDelete)
		apiRouter.HandleFunc("/books/search", api.BooksSearch).Methods(http.MethodGet)

		// Access API
		apiRouter.HandleFunc("/access", api.AccessCreate).Methods(http.MethodPost)
		apiRouter.HandleFunc("/books/{id:[0-9]+}/stats", api.BookStats).Methods(http.MethodGet)
	}

	return r
}
