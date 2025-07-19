package film

import (
	"encoding/json"
	"fmt"
	"log"
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
	w.Header().Set("Content-Type", "application/json")

	films, err := h.service.GetFilms(r.Context())

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(films)
}

func (h *Handler) GetFilmByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

func (h *Handler) SearchFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	title := r.URL.Query().Get("title")

	if title == "" {
		log.Println("no title")
		http.Error(w, "Missing 'title' query parameter", http.StatusBadRequest)
		return
	}

	films, err := h.service.SearchByTitle(r.Context(), title)

	if err != nil {
		log.Println("Scan error:", err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(films)
}
