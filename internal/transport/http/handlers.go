// GET /ui/books/{id}
func (h *Handler) uiBookDetailGET(w http.ResponseWriter, r *http.Request) {
	id := mustUint64(mux.Vars(r)["id"])

	b, err := h.books.Get(r.Context(), id)
	if err != nil {
		h.uiError(w, err)
		return
	}

	stats, _ := h.books.StatsByBook(r.Context(), id)

	// ✅ Convertimos a map[string]int para el template
	statsStr := map[string]int{
		"APERTURA": 0,
		"LECTURA":  0,
		"DESCARGA": 0,
	}
	for k, v := range stats {
		statsStr[string(k)] = v
	}

	data := h.viewBase("Detalle del libro", "book_detail", true)
	data["Book"] = bookToDTO(b)
	data["Stats"] = statsStr                        // ✅ ahora es map[string]int
	data["AccessTypes"] = domain.AllowedAccessTypes // select del form

	h.r.Render(w, "layout", data)
}
