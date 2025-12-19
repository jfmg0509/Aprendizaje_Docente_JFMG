package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

// Renderer renderiza plantillas HTML.
type Renderer struct {
	tpl *template.Template
}

// NewRenderer carga TODOS los templates de web/templates/*.html
// y registra funciones útiles (como tomorrow).
func NewRenderer() (*Renderer, error) {
	funcs := template.FuncMap{
		"tomorrow": func() string {
			return time.Now().Add(24 * time.Hour).Format("2006-01-02")
		},
	}

	pattern := filepath.Join("web", "templates", "*.html")

	tpl, err := template.New("base").Funcs(funcs).ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	return &Renderer{tpl: tpl}, nil
}

// Render ejecuta el template "layout".
// Tu layout debe invocar {{template "content" .}}
func (r *Renderer) Render(w http.ResponseWriter, _ string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.tpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
