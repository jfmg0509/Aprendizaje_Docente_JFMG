package http

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// TemplateRenderer se encarga de cargar y renderizar los templates HTML
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer carga todos los archivos .html de web/templates
func NewTemplateRenderer() (*TemplateRenderer, error) {
	tpl, err := template.ParseGlob(filepath.Join("web", "templates", "*.html"))
	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{
		templates: tpl,
	}, nil
}

// Render renderiza un template específico usando layout + content
func (r *TemplateRenderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Ejecuta el template solicitado
	err := r.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
