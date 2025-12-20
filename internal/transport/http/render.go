package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// Renderer carga templates y ejecuta SIEMPRE "layout.html".
type Renderer struct {
	tpl *template.Template
}

// NewRenderer carga todos los *.html dentro del directorio (ej: web/templates).
func NewRenderer(templatesDir string) (*Renderer, error) {
	pattern := filepath.Join(templatesDir, "*.html")

	tpl, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("parse templates (%s): %w", pattern, err)
	}

	return &Renderer{tpl: tpl}, nil
}

// Render ejecuta un template por nombre (en nuestro caso siempre "layout.html").
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.tpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
