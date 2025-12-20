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
Handler
- Contiene TODOS los handlers UI + API
- NO se repite en ningún otro archivo
*/
type Handler struct {
	userSvc *usecase.UserService
	bookSvc *usecase.BookService
	render  *Renderer
}

// Constructor ÚNICO
func NewHandler(
	userSvc *usecase.UserService,
	bookSvc *usecase.BookService,
	render *Renderer,
) *Handler {
	return &Handler{
		userSvc: userSvc,
		bookSvc: bookSvc,
		render:  render,
	}
}

//
// ========================
// UI (HTML)
// ========================
//

// GET /
func (h *Handler) uiHome(w http.ResponseWriter, r *http.Request) {
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	h.render.Render(w, "home.html", map[string]any{
		"Title":      "Inicio",
		"FooterDate": tomorrow,
	})
}

// GET /ui/users
func (h *Handler) uiUsersGET(w http.ResponseWriter, r *http.Request) {
	users, err := h.userSvc.List(r.Context())
	if err != nil {
		h.uiError(w, err)
		return
	}

	h.render.Render(w, "users.html", map[string]any{
		"Title": "Usuarios",
		"Users": usersToDTO(users),
		"Roles": domain.AllowedRoles,
	})
}

// POST /ui/users
func (h *Handler) uiUsersPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.uiError(w, err)
		return
	}

	_, err := h.userSvc.Create(
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
	books, err := h.bookSvc.List(r.Context())
	if err != nil {
		h.uiError(w, err)
		return
	}

	h.render.Render(w, "books.html", map[string]any{
		"Title": "Libros",
		"Books": booksToDTO(books),
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

	_, err := h.bookSvc.Create(
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
	filter := domain.BookFilter{
		Q:        r.URL.Query().Get("q"),
		Author:   r.URL.Query().Get("author"),
		Category: r.URL.Query().Get("category"),
	}

	books, err := h.bookSvc.Search(r.Context(), filter)
	if err != nil {
		h.uiError(w, err)
		return
	}

	h.render.Render(w, "book_search.html", map[string]any{
		"Title":    "Buscar",
		"Books":    booksToDTO(books),
		"Q":        filter.Q,
		"Author":   filter.Author,
		"Category": filter.Category,
	})
}

// GET /ui/books/{id}
func (h *Handler) uiBookDetailGET(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	book, err := h.bookSvc.Get(r.Context(), id)
	if err != nil {
		h.uiError(w, err)
		return
	}

	stats, _ := h.bookSvc.StatsByBook(r.Context(), id)

	h.render.Render(w, "book_detail.html", map[string]any{
		"Title":       "Detalle del libro",
		"Book":        bookToDTO(book),
		"Stats":       stats,
		"AccessTypes": domain.AllowedAccessTypes,
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

	if err := h.bookSvc.RecordAccess(r.Context(), userID, bookID, accessType); err != nil {
		h.uiError(w, err)
		return
	}

	http.Redirect(w, r, "/ui/books/"+strconv.FormatUint(bookID, 10), http.StatusSeeOther)
}

//
// ========================
// API (JSON)
// ========================
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

	user, err := h.userSvc.Create(r.Context(), in.Name, in.Email, in.Role)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, userToDTO(user))
}

// GET /api/users
func (h *Handler) apiListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userSvc.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, usersToDTO(users))
}

// GET /api/books
func (h *Handler) apiListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookSvc.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, booksToDTO(books))
}

// GET /api/books/search
func (h *Handler) apiSearchBooks(w http.ResponseWriter, r *http.Request) {
	filter := domain.BookFilter{
		Q:        r.URL.Query().Get("q"),
		Author:   r.URL.Query().Get("author"),
		Category: r.URL.Query().Get("category"),
	}

	books, err := h.bookSvc.Search(r.Context(), filter)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, booksToDTO(books))
}

// POST /api/access
func (h *Handler) apiRecordAccess(w http.ResponseWriter, r *http.Request) {
	var in struct {
		UserID     uint64            `json:"user_id"`
		BookID     uint64            `json:"book_id"`
		AccessType domain.AccessType `json:"access_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, err)
		return
	}

	if err := h.bookSvc.RecordAccess(r.Context(), in.UserID, in.BookID, in.AccessType); err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

//
// ========================
// Helpers
// ========================
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

// Render de error UI
func (h *Handler) uiError(w http.ResponseWriter, err error) {
	h.render.Render(w, "error.html", map[string]any{
		"Title": "Error",
		"Error": err.Error(),
	})
}
