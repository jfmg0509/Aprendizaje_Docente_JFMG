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

// ======================================================
// Handler principal
// ======================================================

type Handler struct {
	users *usecase.UserService
	books *usecase.BookService
	r     *Renderer
}

func NewHandler(users *usecase.UserService, books *usecase.BookService, r *Renderer) *Handler {
	return &Handler{
		users: users,
		books: books,
		r:     r,
	}
}

// ======================================================
// API REST (/api/*)
// ======================================================

// ---------- USERS ----------

func (h *Handler) apiCreateUser(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name  string      `json:"name"`
		Email string      `json:"email"`
		Role  domain.Role `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, err)
		return
	}

	u, err := h.users.Create(r.Context(), in.Name, in.Email, in.Role)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, userToDTO(u))
}

func (h *Handler) apiListUsers(w http.ResponseWriter, r *http.Request) {
	list, err := h.users.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, usersToDTO(list))
}

func (h *Handler) apiGetUser(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	u, err := h.users.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, userToDTO(u))
}

func (h *Handler) apiUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	var in struct {
		Name   string      `json:"name"`
		Email  string      `json:"email"`
		Role   domain.Role `json:"role"`
		Active *bool       `json:"active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, err)
		return
	}

	u, err := h.users.Update(r.Context(), id, in.Name, in.Email, in.Role, in.Active)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, userToDTO(u))
}

func (h *Handler) apiDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	if err := h.users.Delete(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ---------- BOOKS ----------

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

	b, err := h.books.Create(
		r.Context(),
		in.Title,
		in.Author,
		in.Year,
		in.ISBN,
		in.Category,
		in.Tags,
		in.Description,
	)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, bookToDTO(b))
}

func (h *Handler) apiListBooks(w http.ResponseWriter, r *http.Request) {
	list, err := h.books.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, booksToDTO(list))
}

func (h *Handler) apiGetBook(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	b, err := h.books.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, bookToDTO(b))
}

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

func (h *Handler) apiDeleteBook(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	if err := h.books.Delete(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ======================================================
// UI HTML (/ui/*)
// ======================================================

// ---------- HOME ----------

func (h *Handler) uiHome(w http.ResponseWriter, r *http.Request) {
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	h.r.Render(w, "home.html", map[string]any{
		"Title":    "Inicio",
		"Tomorrow": tomorrow,
	})
}

// ---------- USERS ----------

func (h *Handler) uiUsersGET(w http.ResponseWriter, r *http.Request) {
	list, err := h.users.List(r.Context())
	if err != nil {
		h.uiError(w, err)
		return
	}

	h.r.Render(w, "users.html", map[string]any{
		"Title": "Usuarios",
		"Users": usersToDTO(list),
		"Roles": domain.AllowedRoles,
	})
}

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

// ---------- BOOKS ----------

func (h *Handler) uiBooksGET(w http.ResponseWriter, r *http.Request) {
	list, err := h.books.List(r.Context())
	if err != nil {
		h.uiError(w, err)
		return
	}

	h.r.Render(w, "books.html", map[string]any{
		"Title": "Libros",
		"Books": booksToDTO(list),
	})
}

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

	h.r.Render(w, "book_search.html", map[string]any{
		"Title":    "Buscar",
		"Books":    booksToDTO(list),
		"Q":        q,
		"Author":   author,
		"Category": category,
	})
}

func (h *Handler) uiBookDetailGET(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	b, err := h.books.Get(r.Context(), id)
	if err != nil {
		h.uiError(w, err)
		return
	}

	stats, _ := h.books.StatsByBook(r.Context(), id)

	h.r.Render(w, "book_detail.html", map[string]any{
		"Title":       "Detalle del libro",
		"Book":        bookToDTO(b),
		"Stats":       stats,
		"AccessTypes": domain.AllowedAccessTypes,
	})
}

// ---------- ERROR ----------

func (h *Handler) uiError(w http.ResponseWriter, err error) {
	h.r.Render(w, "error.html", map[string]any{
		"Title": "Error",
		"Error": err.Error(),
	})
}

// ======================================================
// Helpers
// ======================================================

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
