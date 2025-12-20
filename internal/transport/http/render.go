package http

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type Renderer struct {
	dir   string
	base  *template.Template
	mu    sync.RWMutex
	pages map[string]*template.Template
}

// NewRenderer carga SOLO el layout como base.
// Luego, por cada Render(page), clona el layout y parsea esa página (que define "content").
func NewRenderer(templatesDir string) (*Renderer, error) {
	layoutPath := filepath.Join(templatesDir, "layout.html")

	base, err := template.ParseFiles(layoutPath)
	if err != nil {
		return nil, fmt.Errorf("parse layout (%s): %w", layoutPath, err)
	}

	return &Renderer{
		dir:   templatesDir,
		base:  base,
		pages: map[string]*template.Template{},
	}, nil
}

func (r *Renderer) getPageTemplate(pageFile string) (*template.Template, error) {
	r.mu.RLock()
	if t, ok := r.pages[pageFile]; ok {
		r.mu.RUnlock()
		return t, nil
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	// re-check
	if t, ok := r.pages[pageFile]; ok {
		return t, nil
	}

	clone, err := r.base.Clone()
	if err != nil {
		return nil, fmt.Errorf("clone base: %w", err)
	}

	pagePath := filepath.Join(r.dir, pageFile)
	if _, err := clone.ParseFiles(pagePath); err != nil {
		return nil, fmt.Errorf("parse page (%s): %w", pagePath, err)
	}

	r.pages[pageFile] = clone
	return clone, nil
}

// Render ejecuta "layout" (definido en layout.html) y dentro layout llama a {{template "content" .}}
func (r *Renderer) Render(w http.ResponseWriter, pageFile string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tpl, err := r.getPageTemplate(pageFile)
	if err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
