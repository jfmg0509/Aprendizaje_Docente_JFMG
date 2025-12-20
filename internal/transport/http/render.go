package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// Renderer renderiza HTML cargando SIEMPRE:
// - layout.html
// - la página pedida (ej: users.html)
//
// Esto evita el problema de que todas las páginas definan "content" y se sobreescriban
// cuando parseas todo con ParseGlob.
type Renderer struct {
	dir   string
	mu    sync.RWMutex
	cache map[string]*template.Template // cache por página (users.html, home.html, etc.)
}

// NewRenderer crea el renderer apuntando al directorio web/templates
func NewRenderer(templatesDir string) (*Renderer, error) {
	return &Renderer{
		dir:   templatesDir,
		cache: make(map[string]*template.Template),
	}, nil
}

// Render ejecuta un template por nombre de archivo (ej: "users.html").
// Internamente parsea layout.html + ese archivo y lo cachea.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tpl, err := r.getTemplate(name)
	if err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ejecutamos la página solicitada (ej: users.html).
	// Esa página debe llamar a {{template "layout" .}} internamente.
	if err := tpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template exec error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// getTemplate devuelve el template parseado (layout + page), con cache.
func (r *Renderer) getTemplate(name string) (*template.Template, error) {
	// 1) leer cache
	r.mu.RLock()
	if t, ok := r.cache[name]; ok {
		r.mu.RUnlock()
		return t, nil
	}
	r.mu.RUnlock()

	// 2) construir rutas
	layoutPath := filepath.Join(r.dir, "layout.html")
	pagePath := filepath.Join(r.dir, name)

	// 3) parsear solo lo necesario
	tpl, err := template.ParseFiles(layoutPath, pagePath)
	if err != nil {
		return nil, fmt.Errorf("parse files layout=%s page=%s: %w", layoutPath, pagePath, err)
	}

	// 4) guardar cache
	r.mu.Lock()
	r.cache[name] = tpl
	r.mu.Unlock()

	return tpl, nil
}
