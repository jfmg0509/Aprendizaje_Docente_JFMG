package http

import "net/http"

// renderError muestra errores HTML de forma centralizada
func renderError(w http.ResponseWriter, r *Renderer, err error) {
	r.Render(w, "error.html", map[string]any{
		"Title": "Error",
		"Error": err.Error(),
	})
}
