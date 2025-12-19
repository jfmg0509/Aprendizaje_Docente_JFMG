package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

// UIHandler maneja las páginas HTML (no JSON)
type UIHandler struct {
	renderer *TemplateRenderer
	users    *usecase.UserService
	books    *usecase.BookService
}

// Constructor
func NewUIHandler(
	renderer *TemplateRenderer,
	users *usecase.UserService,
	books *usecase.BookService,
) *UIHandler {
	return &UIHandler{
		renderer: renderer,
		users:    users,
		books:    books,
	}
}

// -------------------- HOME --------------------
func (h *UIHandler) Home(w http.ResponseWriter, r *http.Request) {
	h.renderer.Render(w, "home.html", map[string]any{
		"Title": "Inicio",
	})
}

// -------------------- USERS --------------------
func (h *UIHandler) UsersPage(w http.ResponseWriter, r *http.Request) {
	list, err := h.users.List(r.Context())
	if err != nil {
		h.renderError(w, err)
		return
	}

	h.renderer.Render(w, "users.html", map[string]any{
		"Title": "Usuarios",
		"Users": usersToDTO(list),
		"Roles": domain.AllowedRoles,
	})
}

func (h *UIHandler) UsersCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderError(w, err)
		return
	}

	_, err := h.users.Create(
		r.Context(),
		r.FormValue("name"),
		r.FormValue("email"),
		domain.Role(r.FormValue("role")),
	)
	if err != nil {
		h.renderError(w, err)
		return
	}

	http.Redirect(w, r, "/ui/users", http.StatusSeeOther)
}

// -------------------- BOOKS --------------------
func (h *UIHandler) BooksPage(w http.ResponseWriter, r *http.Request) {
	list, err := h.books.List(r.Context())
	if err != nil {
		h.renderError(w, err)
		return
	}

	h.renderer.Render(w, "books.html", map[string]any{
		"Title": "Libros",
		"Books": booksToDTO(list),
	})
}

func (h *UIHandler) BooksCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.FormValue("year"))
	tags := splitCSV(r.FormValue("tags"))

	_, err := h.books.Create(
		r.Context(),
		r.FormValue("title"),
		r.FormValue("author"),
		year,
		r.FormValue("isbn"),
		r.FormValue("category"),
		tags,
		r.FormValue("description"),
	)
	if err != nil {
		h.renderError(w, err)
		return
	}

	http.Redirect(w, r, "/ui/books", http.StatusSeeOther)
}

// -------------------- SEARCH --------------------
func (h *UIHandler) BookSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	author := r.URL.Query().Get("author")
	category := r.URL.Query().Get("category")

	list, err := h.books.Search(
		r.Context(),
		domain.BookFilter{
			Q:        q,
			Author:   author,
			Category: category,
		},
	)
	if err != nil {
		h.renderError(w, err)
		return
	}

	h.renderer.Render(w, "book_search.html", map[string]any{
		"Title":    "Buscar libros",
		"Books":    booksToDTO(list),
		"Q":        q,
		"Author":   author,
		"Category": category,
	})
}

// -------------------- BOOK DETAIL --------------------
func (h *UIHandler) BookDetail(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseUint(idStr, 10, 64)

	book, err := h.books.Get(r.Context(), id)
	if err != nil {
		h.renderError(w, err)
		return
	}

	stats, _ := h.books.StatsByBook(r.Context(), id)

	h.renderer.Render(w, "book_detail.html", map[string]any{
		"Title":       "Detalle del libro",
		"Book":        bookToDTO(book),
		"Stats":       stats,
		"AccessTypes": domain.AllowedAccessTypes,
	})
}

// -------------------- ERROR --------------------
func (h *UIHandler) renderError(w http.ResponseWriter, err error) {
	h.renderer.Render(w, "error.html", map[string]any{
		"Title": "Error",
		"Error": err.Error(),
	})
}
