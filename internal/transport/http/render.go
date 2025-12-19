package http

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// Renderer carga layout + la página solicitada y luego ejecuta la página.
// Esto evita el error: "no such template content".
type Renderer struct {
	dir string
}

func NewTemplateRenderer() (*Renderer, error) {
	return &Renderer{dir: "web/templates"}, nil
}

// Render carga:
// - web/templates/layout.html
// - web/templates/<name>   (ej: home.html, users.html, books.html, etc)
// y ejecuta el template principal <name>.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	layoutPath := filepath.Join(r.dir, "layout.html")
	pagePath := filepath.Join(r.dir, name)

	tpl, err := template.New("base").ParseFiles(layoutPath, pagePath)
	if err != nil {
		http.Error(w, "template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ejecuta el template principal: "home.html", "users.html", etc.
	if err := tpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template exec error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
