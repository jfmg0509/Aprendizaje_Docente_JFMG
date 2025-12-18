package http

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

type Handlers struct {
	users *usecase.UserService
	books *usecase.BookService
	tpl   *template.Template
}

func NewHandlers(users *usecase.UserService, books *usecase.BookService, tpl *template.Template) *Handlers {
	return &Handlers{users: users, books: books, tpl: tpl}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "time": time.Now().Format(time.RFC3339)})
}

// -------- UI (templates HTML) --------

func (h *Handlers) UIHome(w http.ResponseWriter, r *http.Request) {
	_ = h.tpl.ExecuteTemplate(w, "home.html", map[string]any{"Title": "Sistema de Libros"})
}

func (h *Handlers) UIUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		role := domain.Role(r.FormValue("role"))
		_, _ = h.users.Create(ctx, name, email, role)
		http.Redirect(w, r, "/ui/users", http.StatusSeeOther)
		return
	}

	list, err := h.users.List(ctx)
	if err != nil {
		renderError(w, h.tpl, err)
		return
	}

	_ = h.tpl.ExecuteTemplate(w, "users.html", map[string]any{
		"Title": "Usuarios",
		"Users": list,
		"Roles": domain.AllowedRoles,
	})
}

func (h *Handlers) UIBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		author := r.FormValue("author")
		year, _ := strconv.Atoi(r.FormValue("year"))
		isbn := r.FormValue("isbn")
		category := r.FormValue("category")
		tags := splitCSV(r.FormValue("tags"))
		desc := r.FormValue("description")
		_, _ = h.books.Create(ctx, title, author, year, isbn, category, tags, desc)
		http.Redirect(w, r, "/ui/books", http.StatusSeeOther)
		return
	}

	list, err := h.books.List(ctx)
	if err != nil {
		renderError(w, h.tpl, err)
		return
	}

	_ = h.tpl.ExecuteTemplate(w, "books.html", map[string]any{
		"Title": "Libros",
		"Books": list,
	})
}

func (h *Handlers) UIBookSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query().Get("q")
	author := r.URL.Query().Get("author")
	category := r.URL.Query().Get("category")

	list, err := h.books.Search(ctx, domain.BookFilter{Q: q, Author: author, Category: category})
	if err != nil {
		renderError(w, h.tpl, err)
		return
	}

	_ = h.tpl.ExecuteTemplate(w, "book_search.html", map[string]any{
		"Title": "Búsqueda",
		"Books": list,
		"Q":     q, "Author": author, "Category": category,
	})
}

func (h *Handlers) UIBookDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mustID(w, r)
	if id == 0 {
		return
	}

	b, err := h.books.Get(ctx, id)
	if err != nil {
		renderError(w, h.tpl, err)
		return
	}

	stats, _ := h.books.StatsByBook(ctx, id)

	_ = h.tpl.ExecuteTemplate(w, "book_detail.html", map[string]any{
		"Title":       "Detalle de libro",
		"Book":        b,
		"Stats":       stats,
		"AccessTypes": domain.AllowedAccessTypes,
	})
}

func (h *Handlers) UIAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := strconv.ParseUint(r.FormValue("user_id"), 10, 64)
	bookID, _ := strconv.ParseUint(r.FormValue("book_id"), 10, 64)
	t := domain.AccessType(r.FormValue("access_type"))
	_ = h.books.RecordAccess(ctx, userID, bookID, t)
	http.Redirect(w, r, "/ui/books/"+strconv.FormatUint(bookID, 10), http.StatusSeeOther)
}

// -------- API JSON --------

type apiError struct {
	Error string `json:"error"`
}

func (h *Handlers) APIListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	users, err := h.users.List(ctx)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, usersToDTO(users))
}

func (h *Handlers) APICreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	var in struct {
		Name  string      `json:"name"`
		Email string      `json:"email"`
		Role  domain.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}
	u, err := h.users.Create(ctx, in.Name, in.Email, in.Role)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, userToDTO(u))
}

func (h *Handlers) APIGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	id := mustID(w, r)
	if id == 0 {
		return
	}
	u, err := h.users.Get(ctx, id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, userToDTO(u))
}

func (h *Handlers) APIUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	id := mustID(w, r)
	if id == 0 {
		return
	}

	var in struct {
		Name   *string      `json:"name"`
		Email  *string      `json:"email"`
		Role   *domain.Role `json:"role"`
		Active *bool        `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}
	name := ""
	if in.Name != nil {
		name = *in.Name
	}
	email := ""
	if in.Email != nil {
		email = *in.Email
	}
	role := domain.Role("")
	if in.Role != nil {
		role = *in.Role
	}

	u, err := h.users.Update(ctx, id, name, email, role, in.Active)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, userToDTO(u))
}

func (h *Handlers) APIDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	id := mustID(w, r)
	if id == 0 {
		return
	}
	if err := h.users.Delete(ctx, id); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
}

func (h *Handlers) APIListBooks(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	list, err := h.books.List(ctx)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, booksToDTO(list))
}

func (h *Handlers) APICreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
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
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}
	b, err := h.books.Create(ctx, in.Title, in.Author, in.Year, in.ISBN, in.Category, in.Tags, in.Description)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, bookToDTO(b))
}

func (h *Handlers) APISearchBooks(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	q := r.URL.Query().Get("q")
	author := r.URL.Query().Get("author")
	category := r.URL.Query().Get("category")
	list, err := h.books.Search(ctx, domain.BookFilter{Q: q, Author: author, Category: category})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, booksToDTO(list))
}

func (h *Handlers) APIGetBook(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	id := mustID(w, r)
	if id == 0 {
		return
	}
	b, err := h.books.Get(ctx, id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, bookToDTO(b))
}

func (h *Handlers) APIUpdateBook(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	id := mustID(w, r)
	if id == 0 {
		return
	}
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
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}
	b, err := h.books.Update(ctx, id, usecase.UpdateBookInput{
		Title: in.Title, Author: in.Author, Year: in.Year, ISBN: in.ISBN, Category: in.Category,
		Tags: in.Tags, Description: in.Description, Active: in.Active,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, bookToDTO(b))
}

func (h *Handlers) APIDeleteBook(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	id := mustID(w, r)
	if id == 0 {
		return
	}
	if err := h.books.Delete(ctx, id); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": id})
}

func (h *Handlers) APIRecordAccess(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	var in struct {
		UserID     uint64            `json:"user_id"`
		BookID     uint64            `json:"book_id"`
		AccessType domain.AccessType `json:"access_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}
	if err := h.books.RecordAccess(ctx, in.UserID, in.BookID, in.AccessType); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"queued": true})
}

func (h *Handlers) APIBookStats(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())
	id := mustID(w, r)
	if id == 0 {
		return
	}
	stats, err := h.books.StatsByBook(ctx, id)
	if err != nil {
		writeErr(w, err)
		return
	}
	// maps -> JSON
	writeJSON(w, http.StatusOK, stats)
}

// -------- helpers --------

func withTimeout(ctx context.Context) context.Context {
	c, cancel := context.WithTimeout(ctx, 3*time.Second)

	// cancel automático cuando el contexto termine
	go func() {
		<-c.Done()
		cancel()
	}()

	return c
}

func mustID(w http.ResponseWriter, r *http.Request) uint64 {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid id"})
		return 0
	}
	return id
}

func splitCSV(s string) []string {
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
