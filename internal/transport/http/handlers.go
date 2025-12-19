package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

type UIHandler struct {
	rend  *TemplateRenderer
	users *usecase.UserService
	books *usecase.BookService
}

func NewUIHandler(rend *TemplateRenderer, users *usecase.UserService, books *usecase.BookService) *UIHandler {
	return &UIHandler{rend: rend, users: users, books: books}
}

// HOME
func (h *UIHandler) Home(w http.ResponseWriter, r *http.Request) {
	h.rend.Render(w, "home", "Inicio", nil)
}

// USERS GET
func (h *UIHandler) UsersPage(w http.ResponseWriter, r *http.Request) {
	list, err := h.users.List(r.Context())
	if err != nil {
		h.rend.Render(w, "error", "Error", map[string]any{"Error": err.Error()})
		return
	}

	data := map[string]any{
		"Users": usersToDTO(list),
		"Roles": domain.AllowedRoles,
	}

	h.rend.Render(w, "users", "Usuarios", data)
}

// USERS POST
func (h *UIHandler) UsersCreate(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	role := domain.Role(strings.TrimSpace(r.FormValue("role")))

	_, err := h.users.Create(r.Context(), name, email, role)
	if err != nil {
		h.rend.Render(w, "error", "Error", map[string]any{"Error": err.Error()})
		return
	}

	http.Redirect(w, r, "/ui/users", http.StatusSeeOther)
}

// BOOKS GET
func (h *UIHandler) BooksPage(w http.ResponseWriter, r *http.Request) {
	list, err := h.books.List(r.Context())
	if err != nil {
		h.rend.Render(w, "error", "Error", map[string]any{"Error": err.Error()})
		return
	}

	data := map[string]any{
		"Books": booksToDTO(list),
	}

	h.rend.Render(w, "books", "Libros", data)
}

// BOOKS POST
func (h *UIHandler) BooksCreate(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	title := r.FormValue("title")
	author := r.FormValue("author")
	year, _ := strconv.Atoi(r.FormValue("year"))
	isbn := r.FormValue("isbn")
	category := r.FormValue("category")
	tagsCSV := r.FormValue("tags")
	description := r.FormValue("description")

	var tags []string
	if strings.TrimSpace(tagsCSV) != "" {
		parts := strings.Split(tagsCSV, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				tags = append(tags, p)
			}
		}
	}

	_, err := h.books.Create(r.Context(), title, author, year, isbn, category, tags, description)
	if err != nil {
		h.rend.Render(w, "error", "Error", map[string]any{"Error": err.Error()})
		return
	}

	http.Redirect(w, r, "/ui/books", http.StatusSeeOther)
}

// BOOK DETAIL
func (h *UIHandler) BookDetail(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.ParseUint(idStr, 10, 64)

	b, err := h.books.Get(r.Context(), id)
	if err != nil {
		h.rend.Render(w, "error", "Error", map[string]any{"Error": err.Error()})
		return
	}

	stats, _ := h.books.StatsByBook(r.Context(), id) // si falla, queda vacío
	statsMap := map[string]int{
		"APERTURA": stats[domain.AccessApertura],
		"LECTURA":  stats[domain.AccessLectura],
		"DESCARGA": stats[domain.AccessDescarga],
	}

	data := map[string]any{
		"Book":        bookToDTO(b),
		"AccessTypes": domain.AllowedAccessTypes,
		"Stats":       statsMap,
	}

	h.rend.Render(w, "book_detail", "Detalle del libro", data)
}

// BOOK SEARCH
func (h *UIHandler) BookSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	author := strings.TrimSpace(r.URL.Query().Get("author"))
	category := strings.TrimSpace(r.URL.Query().Get("category"))

	filter := domain.BookFilter{Q: q, Author: author, Category: category}

	list, err := h.books.Search(r.Context(), filter)
	if err != nil {
		h.rend.Render(w, "error", "Error", map[string]any{"Error": err.Error()})
		return
	}

	data := map[string]any{
		"Q":        q,
		"Author":   author,
		"Category": category,
		"Books":    booksToDTO(list),
	}

	h.rend.Render(w, "book_search", "Buscar", data)
}
