package http

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// Renderer carga templates HTML y renderiza páginas.
type Renderer struct {
	t *template.Template
}

// NewRenderer carga todos los templates del folder web/templates.
func NewRenderer(templatesDir string) (*Renderer, error) {
	// Carga *.html
	pattern := filepath.Join(templatesDir, "*.html")
	t, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, err
	}
	return &Renderer{t: t}, nil
}

// Render renderiza un template por nombre (ej: "users.html") con data.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := r.t.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
