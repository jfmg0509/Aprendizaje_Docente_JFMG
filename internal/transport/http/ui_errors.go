package http

import (
    "html/template"
    "net/http"
)

func renderError(w http.ResponseWriter, tpl *template.Template, err error) {
    w.WriteHeader(http.StatusBadRequest)
    _ = tpl.ExecuteTemplate(w, "error.html", map[string]any{
        "Title":"Error",
        "Error": err.Error(),
    })
}
