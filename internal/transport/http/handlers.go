package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

type Handler struct {
	users *usecase.UserService
	books *usecase.BookService
	r     *Renderer
}

func NewHandler(users *usecase.UserService, books *usecase.BookService, r *Renderer) *Handler {
	return &Handler{users: users, books: books, r: r}
}

func (h *Handler) viewBase(title string, showNav bool) map[string]any {
	tomorrow := time.Now().Add(24 * time.Hour).Format("02/01/2006")
	return map[string]any{
		"Title":       title,
		"ShowNav":     showNav,
		"FooterLeft":  "Juan Francisco Morán Gortaire",
		"FooterRight": "PROGRAMACION ORIENTADA A OBJETOS - " + tomorrow,
	}
}

func (h *Handler) uiError(w http.ResponseWriter, err error) {
	data := h.viewBase("Error", true)
	data["Error"] = err.Error()
	h.r.Render(w, "error.html", data)
}

//
// ==============================
// API REST (JSON) - /api/*
// ==============================
//

// POST /api/users
func (h *Handler) apiCreateUser(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name  string      `json:"name"`
		Email string      `json:"email"`
		Role  domain.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, err) // <- ya existe en responses.go
		return
	}

	u, err := h.users.Create(r.Context(), in.Name, in.Email, in.Role)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, userToDTO(u)) // <- ya existe en dto.go
}

// GET /api/users
func (h *Handler) apiListUsers(w http.ResponseWriter, r *http.Request) {
	list, err := h.users.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, usersToDTO(list))
}

// POST /api/books
func (h *Handler) apiCreateBook(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title       string   `json:"title"`
		Author      string   `json:"author"`
		Year        int      `json:"year"`
		ISBN        string   `json:"isbn"`
		Category    string   `json:"category"`
		Tags        []string `json:"tags"`
		Description string   `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, err)
		return
	}

	b, err := h.books.Create(r.Context(), in.Title, in.Author, in.Year, in.ISBN, in.Category, in.Tags, in.Description)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, bookToDTO(b))
}

// GET /api/books
func (h *Handler) apiListBooks(w http.ResponseWriter, r *http.Request) {
	list, err := h.books.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, booksToDTO(list))
}

// GET /api/books/search
func (h *Handler) apiSearchBooks(w http.ResponseWriter, r *http.Request) {
	f := domain.BookFilter{
		Q:        r.URL.Query().Get("q"),
		Author:   r.URL.Query().Get("author"),
		Category: r.URL.Query().Get("category"),
	}
	list, err := h.books.Search(r.Context(), f)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, booksToDTO(list))
}

// GET /api/books/{id}
func (h *Handler) apiGetBook(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])
	b, err := h.books.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, bookToDTO(b))
}

// PATCH /api/books/{id}
func (h *Handler) apiUpdateBook(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	var in struct {
		Title       *string   `json:"title"`
		Author      *string   `json:"author"`
		Year        *int      `json:"year"`
		ISBN        *string   `json:"isbn"`
		Category    *string   `json:"category"`
		Tags        *[]string `json:"tags"`
		Description *string   `json:"description"`
		Active      *bool     `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, err)
		return
	}

	out, err := h.books.Update(r.Context(), id, usecase.UpdateBookInput{
		Title:       in.Title,
		Author:      in.Author,
		Year:        in.Year,
		ISBN:        in.ISBN,
		Category:    in.Category,
		Tags:        in.Tags,
		Description: in.Description,
		Active:      in.Active,
	})
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, bookToDTO(out))
}

// DELETE /api/books/{id}
func (h *Handler) apiDeleteBook(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])
	if err := h.books.Delete(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//
// ==============================
// UI (HTML) - /ui/*
// ==============================
//

// GET /
func (h *Handler) uiHome(w http.ResponseWriter, r *http.Request) {
	data := h.viewBase("Inicio", false)
	h.r.Render(w, "home.html", data)
}

// GET /ui/users
func (h *Handler) uiUsersGET(w http.ResponseWriter, r *http.Request) {
	list, err := h.users.List(r.Context())
	if err != nil {
		h.uiError(w, err)
		return
	}

	data := h.viewBase("Usuarios", true)
	data["Users"] = usersToDTO(list)
	data["Roles"] = domain.AllowedRoles

	h.r.Render(w, "users.html", data)
}

// POST /ui/users
func (h *Handler) uiUsersPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.uiError(w, err)
		return
	}

	_, err := h.users.Create(
		r.Context(),
		r.FormValue("name"),
		r.FormValue("email"),
		domain.Role(r.FormValue("role")),
	)
	if err != nil {
		h.uiError(w, err)
		return
	}

	http.Redirect(w, r, "/ui/users", http.StatusSeeOther)
}

// GET /ui/books
func (h *Handler) uiBooksGET(w http.ResponseWriter, r *http.Request) {
	list, err := h.books.List(r.Context())
	if err != nil {
		h.uiError(w, err)
		return
	}

	data := h.viewBase("Libros", true)
	data["Books"] = booksToDTO(list)

	h.r.Render(w, "books.html", data)
}

// POST /ui/books
func (h *Handler) uiBooksPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.uiError(w, err)
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
		h.uiError(w, err)
		return
	}

	http.Redirect(w, r, "/ui/books", http.StatusSeeOther)
}

// GET /ui/books/search
func (h *Handler) uiBookSearchGET(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	author := r.URL.Query().Get("author")
	category := r.URL.Query().Get("category")

	list, err := h.books.Search(r.Context(), domain.BookFilter{
		Q:        q,
		Author:   author,
		Category: category,
	})
	if err != nil {
		h.uiError(w, err)
		return
	}

	data := h.viewBase("Buscar", true)
	data["Books"] = booksToDTO(list)
	data["Q"] = q
	data["Author"] = author
	data["Category"] = category

	h.r.Render(w, "book_search.html", data)
}

// GET /ui/books/{id}
func (h *Handler) uiBookDetailGET(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	b, err := h.books.Get(r.Context(), id)
	if err != nil {
		h.uiError(w, err)
		return
	}

	data := h.viewBase("Detalle del libro", true)
	data["Book"] = bookToDTO(b)

	// SIN ESTADÍSTICAS
	h.r.Render(w, "book_detail.html", data)
}

//
// ==============================
// Helpers
// ==============================
//

func mustUint64(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}

func splitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
