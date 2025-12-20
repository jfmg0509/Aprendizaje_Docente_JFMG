package http

import "net/http"

// Render HTML de error consistente.
func renderError(r *Renderer, w http.ResponseWriter, title string, err error) {
	r.Render(w, "error.html", map[string]any{
		"Title":       title,
		"Error":       err.Error(),
		"ShowNav":     true,
		"FooterLeft":  "Juan Francisco Morán Gortaire",
		"FooterRight": "", // se rellena en handlers normalmente
	})
}
