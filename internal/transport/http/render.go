package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type Renderer struct {
	dir      string
	layout   *template.Template
	mu       sync.RWMutex
	rendered map[string]*template.Template
}

// NewRenderer
// - Carga SOLO layout.html
// - Las páginas definen "content"
func NewRenderer(templatesDir string) (*Renderer, error) {
	layoutPath := filepath.Join(templatesDir, "layout.html")

	tpl, err := template.ParseFiles(layoutPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing layout.html: %w", err)
	}

	return &Renderer{
		dir:      templatesDir,
		layout:   tpl,
		rendered: make(map[string]*template.Template),
	}, nil
}

// Render
// - Clona layout
// - Parsea la página (users.html, books.html, etc.)
// - Ejecuta template "layout"
func (r *Renderer) Render(w http.ResponseWriter, page string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tpl, err := r.get(page)
	if err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template execute error: "+err.Error(), http.StatusInternalServerError)
	}
}

// get obtiene o construye el template final
func (r *Renderer) get(page string) (*template.Template, error) {
	r.mu.RLock()
	if t, ok := r.rendered[page]; ok {
		r.mu.RUnlock()
		return t, nil
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	// doble chequeo
	if t, ok := r.rendered[page]; ok {
		return t, nil
	}

	clone, err := r.layout.Clone()
	if err != nil {
		return nil, err
	}

	pagePath := filepath.Join(r.dir, page)
	if _, err := clone.ParseFiles(pagePath); err != nil {
		return nil, err
	}

	r.rendered[page] = clone
	return clone, nil
}
