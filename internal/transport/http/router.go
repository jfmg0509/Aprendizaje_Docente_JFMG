package http // Paquete de transporte HTTP (entrada/salida del sistema por web)

import (
	"net/http" // Constantes y tipos HTTP (http.MethodGet, etc.)

	"github.com/gorilla/mux" // Router avanzado (rutas con variables, middlewares, etc.)
)

// Server representa el servidor web del proyecto.
// Contiene el Router (Gorilla Mux) y una referencia a los Handlers (controladores).
type Server struct {
	Router *mux.Router // Router principal: decide qué handler atiende cada ruta
	h      *Handlers   // Handlers: funciones que procesan requests y responden HTML/JSON
}

// NewServer construye el servidor y registra middlewares + rutas.
// Recibe los Handlers ya configurados (con servicios y templates) para conectarlos a las URLs.
func NewServer(h *Handlers) *Server {

	// Crea un router nuevo de Gorilla Mux
	r := mux.NewRouter()

	// Crea la estructura Server con su router y la referencia a handlers
	s := &Server{Router: r, h: h}

	// -----------------------
	// Middlewares globales
	// -----------------------

	// requestIDMiddleware:
	// Genera o inyecta un ID único por request para trazabilidad (logs y debugging).
	r.Use(requestIDMiddleware)

	// loggingMiddleware:
	// Registra en logs info de la request (método, ruta, status, duración, etc.).
	r.Use(loggingMiddleware)

	// methodOverrideMiddleware:
	// Permite que formularios HTML (que solo soportan GET/POST) simulen PUT/DELETE
	// usando un campo hidden o header (útil para la UI).
	r.Use(methodOverrideMiddleware)

	// -----------------------
	// Health check
	// -----------------------

	// Endpoint simple para verificar que el servidor está vivo.
	// GET /health -> responde OK (normalmente 200)
	r.HandleFunc("/health", h.Health).Methods(http.MethodGet)

	// -----------------------
	// UI (Frontend HTML)
	// -----------------------

	// Home de la UI
	r.HandleFunc("/", h.UIHome).Methods(http.MethodGet)

	// Lista usuarios (GET) y crea usuario desde formulario (POST)
	r.HandleFunc("/ui/users", h.UIUsers).Methods(http.MethodGet, http.MethodPost)

	// Lista libros (GET) y crea libro desde formulario (POST)
	r.HandleFunc("/ui/books", h.UIBooks).Methods(http.MethodGet, http.MethodPost)

	// Búsqueda de libros desde UI (GET con query params)
	r.HandleFunc("/ui/books/search", h.UIBookSearch).Methods(http.MethodGet)

	// Detalle de libro por ID (solo números por el regex [0-9]+)
	r.HandleFunc("/ui/books/{id:[0-9]+}", h.UIBookDetail).Methods(http.MethodGet)

	// Registrar un acceso desde UI (POST)
	// Esto normalmente dispara el registro en la cola concurrente (goroutines/canales).
	r.HandleFunc("/ui/access", h.UIAccess).Methods(http.MethodPost)

	// -----------------------
	// API Users (JSON / REST)
	// -----------------------

	// GET /api/users -> lista todos los usuarios (JSON)
	r.HandleFunc("/api/users", h.APIListUsers).Methods(http.MethodGet)

	// POST /api/users -> crea un usuario (JSON)
	r.HandleFunc("/api/users", h.APICreateUser).Methods(http.MethodPost)

	// GET /api/users/{id} -> obtiene un usuario por ID
	r.HandleFunc("/api/users/{id:[0-9]+}", h.APIGetUser).Methods(http.MethodGet)

	// PUT /api/users/{id} -> actualiza un usuario por ID
	r.HandleFunc("/api/users/{id:[0-9]+}", h.APIUpdateUser).Methods(http.MethodPut)

	// DELETE /api/users/{id} -> elimina un usuario por ID
	r.HandleFunc("/api/users/{id:[0-9]+}", h.APIDeleteUser).Methods(http.MethodDelete)

	// -----------------------
	// API Books (JSON / REST)
	// -----------------------

	// GET /api/books -> lista libros
	r.HandleFunc("/api/books", h.APIListBooks).Methods(http.MethodGet)

	// POST /api/books -> crea libro
	r.HandleFunc("/api/books", h.APICreateBook).Methods(http.MethodPost)

	// GET /api/books/search -> búsqueda de libros (por query params)
	r.HandleFunc("/api/books/search", h.APISearchBooks).Methods(http.MethodGet)

	// GET /api/books/{id} -> obtiene libro por ID
	r.HandleFunc("/api/books/{id:[0-9]+}", h.APIGetBook).Methods(http.MethodGet)

	// PUT /api/books/{id} -> actualiza libro
	r.HandleFunc("/api/books/{id:[0-9]+}", h.APIUpdateBook).Methods(http.MethodPut)

	// DELETE /api/books/{id} -> elimina libro
	r.HandleFunc("/api/books/{id:[0-9]+}", h.APIDeleteBook).Methods(http.MethodDelete)

	// GET /api/books/{id}/stats -> devuelve estadísticas (ej: vistas/descargas por tipo)
	r.HandleFunc("/api/books/{id:[0-9]+}/stats", h.APIBookStats).Methods(http.MethodGet)

	// -----------------------
	// API Access (evento de acceso)
	// -----------------------

	// POST /api/access -> registra un acceso (ej: view/download)
	// Normalmente se encola y lo procesan goroutines para no bloquear el request.
	r.HandleFunc("/api/access", h.APIRecordAccess).Methods(http.MethodPost)

	// Devuelve el Server con router listo para usarse en main.go
	return s
}
