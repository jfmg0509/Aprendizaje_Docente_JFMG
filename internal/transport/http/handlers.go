package http // Paquete de transporte HTTP: recibe requests, valida datos y devuelve respuestas (HTML/JSON)

import (
	"context"       // Contextos para control de timeout/cancelación de operaciones
	"encoding/json" // Codificar/decodificar JSON (API REST)
	"html/template" // Templates HTML (UI)
	"net/http"      // Tipos HTTP, ResponseWriter, Request, métodos, status codes
	"strconv"       // Conversión de string a int/uint y viceversa
	"strings"       // Manipulación de cadenas (split/trim)
	"time"          // Timeouts y fecha/hora

	"github.com/gorilla/mux"                                           // Router Gorilla Mux: variables en rutas, middlewares, etc.
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"  // Dominio: Roles, AccessTypes, filtros, etc.
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase" // Casos de uso: servicios de negocio
)

// Handlers agrupa todos los endpoints del sistema (UI + API).
// Contiene referencias a los servicios (usecase) y a los templates HTML.
type Handlers struct {
	users *usecase.UserService // Servicio de usuarios (lógica de negocio)
	books *usecase.BookService // Servicio de libros (lógica de negocio)
	tpl   *template.Template   // Motor de templates HTML para la UI
}

// NewHandlers es el "constructor" del struct Handlers.
// Inyecta dependencias: servicios y templates, para no acoplar handlers a DB directamente.
func NewHandlers(users *usecase.UserService, books *usecase.BookService, tpl *template.Template) *Handlers {
	return &Handlers{users: users, books: books, tpl: tpl}
}

// Health responde si el servidor está vivo (endpoint /health).
// Útil para monitoreo, pruebas y verificación de despliegue.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	// Responde JSON con status y timestamp en formato RFC3339.
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// -------------------- UI (templates HTML) --------------------

// UIHome renderiza la página principal (GET /).
func (h *Handlers) UIHome(w http.ResponseWriter, r *http.Request) {
	// Calculamos la fecha de mañana
	tomorrow := time.Now().
		Add(24 * time.Hour).
		Format("02/01/2006") // formato DD/MM/YYYY

	// Enviamos datos al template
	err := h.tpl.ExecuteTemplate(w, "home.html", map[string]any{
		"Title":        "Inicio",
		"TomorrowDate": tomorrow,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UIUsers maneja la pantalla de usuarios (GET lista / POST crea).
// Ruta: /ui/users
func (h *Handlers) UIUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Contexto de la request (cancelación si el cliente corta)

	// Si viene POST, significa que el formulario envió datos para crear usuario.
	if r.Method == http.MethodPost {
		// Extrae campos del formulario HTML (name, email, role)
		name := r.FormValue("name")
		email := r.FormValue("email")
		role := domain.Role(r.FormValue("role")) // Convierte a tipo Role del dominio

		// Crea usuario usando la lógica del servicio (usecase).
		// Aquí se ignoran errores por simplicidad UI (en API sí se reporta).
		_, _ = h.users.Create(ctx, name, email, role)

		// Redirige para evitar reenvío del formulario (patrón PRG: Post-Redirect-Get).
		http.Redirect(w, r, "/ui/users", http.StatusSeeOther)
		return
	}

	// Si no es POST, es GET: listamos usuarios.
	list, err := h.users.List(ctx)
	if err != nil {
		// Renderiza una página HTML de error usando templates.
		renderError(w, h.tpl, err)
		return
	}

	// Renderiza "users.html" con la lista de usuarios y roles permitidos.
	_ = h.tpl.ExecuteTemplate(w, "users.html", map[string]any{
		"Title": "Usuarios",
		"Users": list,
		"Roles": domain.AllowedRoles, // Lista de roles válidos (array/slice del dominio)
	})
}

// UIBooks maneja la pantalla de libros (GET lista / POST crea).
// Ruta: /ui/books
func (h *Handlers) UIBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// POST => crear libro desde formulario.
	if r.Method == http.MethodPost {
		// Campos del formulario
		title := r.FormValue("title")
		author := r.FormValue("author")

		// Convierte el año a int; si falla queda 0 (se podría validar más).
		year, _ := strconv.Atoi(r.FormValue("year"))

		isbn := r.FormValue("isbn")
		category := r.FormValue("category")

		// Tags vienen como string separado por comas: "go, poo, backend"
		tags := splitCSV(r.FormValue("tags"))

		desc := r.FormValue("description")

		// Crea el libro mediante el servicio.
		_, _ = h.books.Create(ctx, title, author, year, isbn, category, tags, desc)

		// PRG: redirigir para evitar doble envío.
		http.Redirect(w, r, "/ui/books", http.StatusSeeOther)
		return
	}

	// GET => listar libros.
	list, err := h.books.List(ctx)
	if err != nil {
		renderError(w, h.tpl, err)
		return
	}

	// Renderiza la vista con la lista.
	_ = h.tpl.ExecuteTemplate(w, "books.html", map[string]any{
		"Title": "Libros",
		"Books": list,
	})
}

// UIBookSearch muestra resultados de búsqueda (GET).
// Ruta: /ui/books/search?q=...&author=...&category=...
func (h *Handlers) UIBookSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtiene query params
	q := r.URL.Query().Get("q")
	author := r.URL.Query().Get("author")
	category := r.URL.Query().Get("category")

	// Ejecuta búsqueda con filtro del dominio.
	list, err := h.books.Search(ctx, domain.BookFilter{
		Q:        q,
		Author:   author,
		Category: category,
	})
	if err != nil {
		renderError(w, h.tpl, err)
		return
	}

	// Renderiza resultados y mantiene filtros en la vista.
	_ = h.tpl.ExecuteTemplate(w, "book_search.html", map[string]any{
		"Title":    "Búsqueda",
		"Books":    list,
		"Q":        q,
		"Author":   author,
		"Category": category,
	})
}

// UIBookDetail muestra el detalle de un libro y sus estadísticas.
// Ruta: /ui/books/{id}
func (h *Handlers) UIBookDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extrae y valida el ID de la ruta usando mux vars.
	id := mustID(w, r)
	if id == 0 {
		// mustID ya respondió el error 400 si aplica.
		return
	}

	// Obtiene el libro por ID.
	b, err := h.books.Get(ctx, id)
	if err != nil {
		renderError(w, h.tpl, err)
		return
	}

	// Obtiene estadísticas (por ejemplo, accesos por tipo).
	// Se ignora error para no bloquear UI (podrías manejarlo si deseas).
	stats, _ := h.books.StatsByBook(ctx, id)

	// Renderiza la vista con libro, stats y tipos de acceso disponibles.
	_ = h.tpl.ExecuteTemplate(w, "book_detail.html", map[string]any{
		"Title":       "Detalle de libro",
		"Book":        b,
		"Stats":       stats,                     // map/slice con estadísticas
		"AccessTypes": domain.AllowedAccessTypes, // lista de tipos permitidos (domain)
	})
}

// UIAccess registra un acceso (view/download) desde un formulario HTML.
// Ruta: /ui/access (POST)
func (h *Handlers) UIAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Obtiene IDs desde formulario (user_id, book_id).
	userID, _ := strconv.ParseUint(r.FormValue("user_id"), 10, 64)
	bookID, _ := strconv.ParseUint(r.FormValue("book_id"), 10, 64)

	// Tipo de acceso (view/download/etc.) como tipo del dominio.
	t := domain.AccessType(r.FormValue("access_type"))

	// Registra el acceso; por diseño puede encolarse (concurrencia).
	_ = h.books.RecordAccess(ctx, userID, bookID, t)

	// Redirige al detalle del libro.
	http.Redirect(w, r, "/ui/books/"+strconv.FormatUint(bookID, 10), http.StatusSeeOther)
}

// -------------------- API (JSON / REST) --------------------

// apiError define la forma estándar de error JSON.
type apiError struct {
	Error string `json:"error"` // Campo "error" en JSON
}

// APIListUsers lista usuarios en JSON.
// GET /api/users
func (h *Handlers) APIListUsers(w http.ResponseWriter, r *http.Request) {
	// Aplica timeout al contexto para que no se quede colgado.
	ctx := withTimeout(r.Context())

	// Llama al servicio (usecase).
	users, err := h.users.List(ctx)
	if err != nil {
		writeErr(w, err) // Respuesta uniforme de errores
		return
	}

	// Convierte entidad de dominio a DTO antes de responder JSON.
	writeJSON(w, http.StatusOK, usersToDTO(users))
}

// APICreateUser crea usuario desde JSON.
// POST /api/users
func (h *Handlers) APICreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())

	// Estructura de entrada esperada (JSON).
	var in struct {
		Name  string      `json:"name"`
		Email string      `json:"email"`
		Role  domain.Role `json:"role"`
	}

	// Decodifica el JSON del body a la estructura.
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}

	// Crea usuario.
	u, err := h.users.Create(ctx, in.Name, in.Email, in.Role)
	if err != nil {
		writeErr(w, err)
		return
	}

	// Devuelve el usuario creado en formato DTO.
	writeJSON(w, http.StatusCreated, userToDTO(u))
}

// APIGetUser obtiene un usuario por ID.
// GET /api/users/{id}
func (h *Handlers) APIGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())

	// Obtiene ID y valida.
	id := mustID(w, r)
	if id == 0 {
		return
	}

	// Busca usuario por ID.
	u, err := h.users.Get(ctx, id)
	if err != nil {
		writeErr(w, err)
		return
	}

	// Responde con DTO.
	writeJSON(w, http.StatusOK, userToDTO(u))
}

// APIUpdateUser actualiza usuario por ID.
// PUT /api/users/{id}
func (h *Handlers) APIUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())

	// ID desde ruta
	id := mustID(w, r)
	if id == 0 {
		return
	}

	// Se usan punteros para permitir "campos opcionales" en JSON.
	// Si el campo es nil => no se actualiza.
	var in struct {
		Name   *string      `json:"name"`
		Email  *string      `json:"email"`
		Role   *domain.Role `json:"role"`
		Active *bool        `json:"active"`
	}

	// Decodifica JSON
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}

	// Normaliza campos opcionales: si nil => "", si existe => valor.
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

	// Llama al servicio con los valores preparados.
	u, err := h.users.Update(ctx, id, name, email, role, in.Active)
	if err != nil {
		writeErr(w, err)
		return
	}

	// Responde con DTO actualizado.
	writeJSON(w, http.StatusOK, userToDTO(u))
}

// APIDeleteUser elimina un usuario.
// DELETE /api/users/{id}
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

// APIListBooks lista libros.
// GET /api/books
func (h *Handlers) APIListBooks(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())

	list, err := h.books.List(ctx)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, booksToDTO(list))
}

// APICreateBook crea un libro.
// POST /api/books
func (h *Handlers) APICreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())

	// Estructura esperada del JSON de entrada.
	var in struct {
		Title       string   `json:"title"`
		Author      string   `json:"author"`
		Year        int      `json:"year"`
		ISBN        string   `json:"isbn"`
		Category    string   `json:"category"`
		Tags        []string `json:"tags"`
		Description string   `json:"description"`
	}

	// Decodifica JSON.
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}

	// Crea el libro.
	b, err := h.books.Create(ctx, in.Title, in.Author, in.Year, in.ISBN, in.Category, in.Tags, in.Description)
	if err != nil {
		writeErr(w, err)
		return
	}

	// Devuelve el libro creado.
	writeJSON(w, http.StatusCreated, bookToDTO(b))
}

// APISearchBooks busca libros por filtros.
// GET /api/books/search?q=...&author=...&category=...
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

// APIGetBook obtiene un libro por ID.
// GET /api/books/{id}
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

// APIUpdateBook actualiza un libro por ID.
// PUT /api/books/{id}
func (h *Handlers) APIUpdateBook(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())

	id := mustID(w, r)
	if id == 0 {
		return
	}

	// Campos opcionales con punteros: si nil => no se actualiza.
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

	// Llama al servicio usando un input struct (más limpio y extensible).
	b, err := h.books.Update(ctx, id, usecase.UpdateBookInput{
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

	writeJSON(w, http.StatusOK, bookToDTO(b))
}

// APIDeleteBook elimina libro por ID.
// DELETE /api/books/{id}
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

// APIRecordAccess registra un acceso (evento) en el sistema.
// POST /api/access
func (h *Handlers) APIRecordAccess(w http.ResponseWriter, r *http.Request) {
	ctx := withTimeout(r.Context())

	// Entrada esperada.
	var in struct {
		UserID     uint64            `json:"user_id"`
		BookID     uint64            `json:"book_id"`
		AccessType domain.AccessType `json:"access_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}

	// Registra acceso (posiblemente encolado en goroutines/channels).
	if err := h.books.RecordAccess(ctx, in.UserID, in.BookID, in.AccessType); err != nil {
		writeErr(w, err)
		return
	}

	// 202 Accepted indica que fue aceptado para procesamiento (ideal si es asíncrono).
	writeJSON(w, http.StatusAccepted, map[string]any{"queued": true})
}

// APIBookStats devuelve estadísticas de un libro.
// GET /api/books/{id}/stats
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

	// stats suele ser un map => JSON directo.
	writeJSON(w, http.StatusOK, stats)
}

// -------------------- Helpers --------------------

// withTimeout crea un contexto con límite de tiempo.
// Esto evita que una request quede colgada indefinidamente (por DB lenta, etc.).
func withTimeout(ctx context.Context) context.Context {
	// Crea un contexto hijo que expira en 3 segundos.
	c, cancel := context.WithTimeout(ctx, 3*time.Second)

	// Como esta función no devuelve cancel(), hacemos un "auto-cancel":
	// Cuando el contexto termina (timeout o cancelación), ejecutamos cancel()
	// para liberar recursos internos (timer).
	go func() {
		<-c.Done() // espera a que el contexto finalice
		cancel()   // libera recursos
	}()

	return c // Retorna el contexto con timeout para ser usado en servicios/repos
}

// mustID lee el {id} de la URL y lo valida (debe ser uint64 > 0).
// Si no es válido, responde 400 con JSON de error y retorna 0.
func mustID(w http.ResponseWriter, r *http.Request) uint64 {
	// mux.Vars obtiene variables de la ruta, ej: /api/books/{id}
	idStr := mux.Vars(r)["id"]

	// Convierte a número
	id, err := strconv.ParseUint(idStr, 10, 64)

	// Validación básica
	if err != nil || id == 0 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid id"})
		return 0
	}

	return id
}

// splitCSV convierte "a, b, c" => []string{"a","b","c"} limpiando espacios y vacíos.
// Útil para tags ingresados en formularios.
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")       // separa por coma
	out := make([]string, 0, len(parts)) // slice con capacidad inicial (eficiencia)

	for _, p := range parts { // recorre cada parte
		p = strings.TrimSpace(p) // elimina espacios al inicio/fin
		if p != "" {             // ignora elementos vacíos
			out = append(out, p) // agrega al slice
		}
	}

	return out // devuelve lista final de tags
}
