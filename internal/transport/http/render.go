package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

// TemplateRenderer carga y renderiza templates HTML.
type TemplateRenderer struct {
	tpl *template.Template
}

// NewTemplateRenderer parsea TODOS los templates de web/templates/*.html
func NewTemplateRenderer() (*TemplateRenderer, error) {
	pattern := filepath.Join("web", "templates", "*.html")

	tpl, err := template.New("").ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	return &TemplateRenderer{tpl: tpl}, nil
}

// ViewData es el “sobre” común que viaja a todos los templates.
type ViewData struct {
	Title    string // título del <title> y/o cabecera
	View     string // nombre del template a renderizar: "home", "users", etc.
	Tomorrow string // fecha de mañana para el footer del Home

	// Data es donde metemos tus datos reales (Users, Books, etc.)
	Data any
}

// render ejecuta el layout y dentro del layout se dibuja {{template .View .}}
func (r *TemplateRenderer) Render(w http.ResponseWriter, view string, title string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vd := ViewData{
		Title: title,
		View:  view,
		Data:  data,
	}

	// Si es home, enviamos la fecha de mañana
	if view == "home" {
		vd.Tomorrow = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}

	// Ejecuta el layout principal
	if err := r.tpl.ExecuteTemplate(w, "layout", vd); err != nil {
		http.Error(w, "template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
