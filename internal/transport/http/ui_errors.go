package http

import "net/http"

// renderError: helper para mostrar un error en una página HTML estándar.
// Esto evita repetir el mismo map[string]any en varios handlers.
func renderError(r *Renderer, w http.ResponseWriter, err error) {
	// Puedes cambiar el status si quieres. Por simplicidad dejamos 200,
	// pero podrías poner 400/500 según el error.
	// w.WriteHeader(http.StatusInternalServerError)

	r.Render(w, "error.html", map[string]any{
		"Title": "Error",
		"Error": err.Error(),
	})
}
