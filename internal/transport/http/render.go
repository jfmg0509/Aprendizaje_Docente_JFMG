package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	dir string
}

// NewRenderer crea el renderer apuntando a la carpeta de templates.
// Ej: NewRenderer("web/templates")
func NewRenderer(templatesDir string) *Renderer {
	return &Renderer{dir: templatesDir}
}

// Render carga layout.html + la página solicitada (ej: home.html) y ejecuta "layout".
// Así evitamos que varios {{define "content"}} se pisen entre sí.
func (r *Renderer) Render(w http.ResponseWriter, page string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	layoutPath := filepath.Join(r.dir, "layout.html")
	pagePath := filepath.Join(r.dir, page)

	tpl, err := template.ParseFiles(layoutPath, pagePath)
	if err != nil {
		http.Error(w, "templates: "+fmt.Errorf("parse files: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	// Ejecuta el layout (que incluye {{template "content" .}})
	if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
