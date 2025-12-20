package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// Renderer carga y renderiza templates HTML.
// Convención:
//   - layout.html define el layout general.
//   - Cada página (home.html, users.html, etc.) define su propio template name igual al archivo,
//     y dentro llama a {{template "layout" .}} o directamente incluye el contenido.
//
// En nuestro caso: Render() ejecuta el template por nombre.
type Renderer struct {
	tpl *template.Template
}

// NewRenderer carga todos los templates *.html dentro del directorio indicado.
// Ejemplo: NewRenderer("web/templates")
func NewRenderer(templatesDir string) (*Renderer, error) {
	pattern := filepath.Join(templatesDir, "*.html")

	tpl, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("parse templates (%s): %w", pattern, err)
	}

	return &Renderer{tpl: tpl}, nil
}

// Render ejecuta un template por nombre (ej: "home.html", "users.html").
// data es el map/struct que mandas desde los handlers.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.tpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
