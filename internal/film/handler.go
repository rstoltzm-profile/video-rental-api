package film

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appication/json")

	films, err := h.service.GetFilms(r.Context())

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(films)
}
