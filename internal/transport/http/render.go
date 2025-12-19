package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

// Renderer se encarga de cargar templates y renderizar páginas HTML
type Renderer struct {
	baseDir string             // carpeta donde están los templates (web/templates)
	layout  *template.Template // template base (layout)
}

// NewTemplateRenderer crea un Renderer listo para usar.
//
// Estructura esperada:
// web/templates/layout.html   (define "layout" y usa {{template "content" .}})
// web/templates/home.html     (define "content")
// web/templates/users.html    (define "content")
// web/templates/books.html    (define "content")
// web/templates/book_search.html (define "content")
// web/templates/book_detail.html (define "content")
// web/templates/error.html    (define "content")
func NewTemplateRenderer() (*Renderer, error) {
	baseDir := filepath.Join("web", "templates")

	funcs := template.FuncMap{
		// Fecha de mañana (por si la necesitas en footer)
		"tomorrow": func() string {
			return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		},
		// Fecha de hoy
		"today": func() string {
			return time.Now().Format("2006-01-02")
		},
	}

	// Cargamos SOLO el layout como base
	layoutPath := filepath.Join(baseDir, "layout.html")
	tpl, err := template.New("layout.html").Funcs(funcs).ParseFiles(layoutPath)
	if err != nil {
		return nil, fmt.Errorf("parse layout: %w", err)
	}

	return &Renderer{
		baseDir: baseDir,
		layout:  tpl,
	}, nil
}

// Render renderiza una vista (ej: "users.html") dentro del layout.
// Internamente:
// 1) Clona el layout
// 2) Parse el archivo de vista (que define "content")
// 3) Ejecuta el template "layout"
func (r *Renderer) Render(w http.ResponseWriter, view string, data any) {
	// Clonar el layout para no mezclar estados entre requests
	tpl, err := r.layout.Clone()
	if err != nil {
		http.Error(w, "template clone error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parsear la vista sobre ese clon
	viewPath := filepath.Join(r.baseDir, view)
	if _, err := tpl.ParseFiles(viewPath); err != nil {
		http.Error(w, "template parse view error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Ejecutar siempre el layout
	if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template execute error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
