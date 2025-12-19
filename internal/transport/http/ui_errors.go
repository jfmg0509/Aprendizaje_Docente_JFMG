package http

import (
	"net/http"
)

// renderError es un helper por si quieres usarlo desde middlewares o handlers.
// (Ahora mismo, el Handler ya tiene uiError, pero esto queda útil.)
func renderError(w http.ResponseWriter, r *Renderer, err error) {
	r.Render(w, "error.html", map[string]any{
		"Title": "Error",
		"Error": err.Error(),
	})
}
