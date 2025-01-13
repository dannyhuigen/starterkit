package handlers

import (
	"net/http"
	"starterkit/internal/templates/layouts"
	"starterkit/internal/templates/pages"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	layoutWithContent := layouts.Mainlayout(pages.Index())
	err := layoutWithContent.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
