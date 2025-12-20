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

/*
Handler centraliza TODOS los endpoints:
- UI (HTML)
- API (JSON)
*/
type Handler struct {
	users *usecase.UserService
	books *usecase.BookService
	r     *Renderer
}

/*
Constructor único del handler
*/
func NewHandler(
	userSvc *usecase.UserService,
	bookSvc *usecase.BookService,
	renderer *Renderer,
) *Handler {
	return &Handler{
		users: userSvc,
		books: bookSvc,
		r:     renderer,
	}
}

//
// =====================================================
// ======================= UI ==========================
// =====================================================
//

// GET /
func (h *Handler) uiHome(w http.ResponseWriter, r *http.Request) {
	h.r.Render(w, "layout.html", map[string]any{
		"Title":    "Inicio",
		"View":     "home",
		"Tomorrow": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
	})
}

// GET /ui/users
func (h *Handler) uiUsersGET(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.List(r.Context())
	if err != nil {
		h.uiError(w, err)
		return
	}

	h.r.Render(w, "layout.html", map[string]any{
		"Title":    "Usuarios",
		"View":     "users",
		"Users":    usersToDTO(users),
		"Roles":    domain.AllowedRoles,
		"Tomorrow": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
	})
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
	books, err := h.books.List(r.Context())
	if err != nil {
		h.uiError(w, err)
		return
	}

	h.r.Render(w, "layout.html", map[string]any{
		"Title":    "Libros",
		"View":     "books",
		"Books":    booksToDTO(books),
		"Tomorrow": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
	})
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

	books, err := h.books.Search(r.Context(), domain.BookFilter{
		Q:        q,
		Author:   author,
		Category: category,
	})
	if err != nil {
		h.uiError(w, err)
		return
	}

	h.r.Render(w, "layout.html", map[string]any{
		"Title":    "Buscar",
		"View":     "book_search",
		"Books":    booksToDTO(books),
		"Q":        q,
		"Author":   author,
		"Category": category,
		"Tomorrow": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
	})
}

// GET /ui/books/{id}
func (h *Handler) uiBookDetailGET(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	book, err := h.books.Get(r.Context(), id)
	if err != nil {
		h.uiError(w, err)
		return
	}

	stats, _ := h.books.StatsByBook(r.Context(), id)

	h.r.Render(w, "layout.html", map[string]any{
		"Title":       "Detalle del libro",
		"View":        "book_detail",
		"Book":        bookToDTO(book),
		"Stats":       stats,
		"AccessTypes": domain.AllowedAccessTypes,
		"Tomorrow":    time.Now().Add(24 * time.Hour).Format("2006-01-02"),
	})
}

// POST /ui/access
func (h *Handler) uiAccessPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.uiError(w, err)
		return
	}

	userID, _ := strconv.ParseUint(r.FormValue("user_id"), 10, 64)
	bookID, _ := strconv.ParseUint(r.FormValue("book_id"), 10, 64)
	accessType := domain.AccessType(r.FormValue("access_type"))

	if err := h.books.RecordAccess(r.Context(), userID, bookID, accessType); err != nil {
		h.uiError(w, err)
		return
	}

	http.Redirect(w, r, "/ui/books/"+strconv.FormatUint(bookID, 10), http.StatusSeeOther)
}

//
// =====================================================
// ======================= API =========================
// =====================================================
//

// POST /api/users
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

// GET /api/users
func (h *Handler) apiListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, usersToDTO(users))
}

// GET /api/users/{id}
func (h *Handler) apiGetUser(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	user, err := h.users.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, userToDTO(user))
}

// DELETE /api/users/{id}
func (h *Handler) apiDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	if err := h.users.Delete(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	book, err := h.books.Create(
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

	writeJSON(w, http.StatusCreated, bookToDTO(book))
}

// GET /api/books/{id}/stats
func (h *Handler) apiStatsByBook(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	stats, err := h.books.StatsByBook(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

//
// =====================================================
// ===================== HELPERS =======================
// =====================================================
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

func (h *Handler) uiError(w http.ResponseWriter, err error) {
	h.r.Render(w, "layout.html", map[string]any{
		"Title":    "Error",
		"View":     "error",
		"Error":    err.Error(),
		"Tomorrow": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
	})
}
