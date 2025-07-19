package film

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func (h *Handler) GetFilmByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "appication/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
	}
	film, err := h.service.GetFilmByID(r.Context(), id)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(film)
}
