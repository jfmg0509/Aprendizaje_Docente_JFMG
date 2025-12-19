package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// Renderer se encarga de cargar y renderizar templates HTML.
// Convención usada:
// - layout.html define: {{define "layout"}} ... {{template "content" .}} ... {{end}}
// - cada página define:
//  1. {{define "NOMBRE_DE_LA_PAGINA"}} {{template "layout" .}} {{end}}
//  2. {{define "content"}} ... contenido ... {{end}}
//
// Importante:
//   - Render() NO ejecuta "layout" directamente, porque "layout" necesita que exista
//     un template "content" ya definido. En su lugar, ejecuta el template de la página.
type Renderer struct {
	tpl *template.Template
	dir string
}

// NewRenderer carga todos los templates *.html dentro del directorio indicado.
// Ejemplo: NewRenderer("web/templates")
func NewRenderer(templatesDir string) (*Renderer, error) {
	pattern := filepath.Join(templatesDir, "*.html")

	tpl, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("parse templates (%s): %w", pattern, err)
	}

	return &Renderer{
		tpl: tpl,
		dir: templatesDir,
	}, nil
}

// Render ejecuta un template por nombre (ej: "users.html", "home.html").
// data es el map/struct que tú mandas desde los handlers.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Ejecutamos el template de la página, NO el layout.
	// Ese template de página debe llamar a {{template "layout" .}}
	if err := r.tpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
