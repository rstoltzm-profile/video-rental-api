package film

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetFilms godoc
// @Summary      List all films
// @Description  Returns a list of all films
// @Tags         films
// @Produce      json
// @Success      200  {array}   film.Film
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /films [get]
func (h *Handler) GetFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	films, err := h.service.GetFilms(r.Context())

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(films)
}

// GetFilmByID godoc
// @Summary      Get a film by ID
// @Description  Returns a single film by its ID
// @Tags         films
// @Produce      json
// @Param        id   path      int  true  "Film ID"
// @Success      200  {object}  film.Film
// @Failure      400  {string}  string "Invalid film ID"
// @Failure      404  {string}  string "Film not found"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /films/{id} [get]
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

// SearchFilm godoc
// @Summary      Search films by title
// @Description  Returns a list of films matching the title query parameter
// @Tags         films
// @Produce      json
// @Param        title  query     string  true  "Film title to search for"
// @Success      200    {array}   film.Film
// @Failure      400    {string}  string "Missing title query parameter"
// @Failure      500    {string}  string "Internal Server Error"
// @Router       /films/search [get]
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

// GetFilmWithActorsAndCategoriesByID godoc
// @Summary      Get film with actors and categories by ID
// @Description  Returns a film along with its actors and categories
// @Tags         films
// @Produce      json
// @Param        id   path      int  true  "Film ID"
// @Success      200  {object}  film.FilmWithActorsCategories
// @Failure      400  {string}  string "Invalid film ID"
// @Failure      404  {string}  string "Film not found"
// @Router       /films/{id}/with-actors-categories [get]
func (h *Handler) GetFilmWithActorsAndCategoriesByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 || parts[3] != "with-actors-categories" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid film ID", http.StatusBadRequest)
		return
	}

	filmWithActors, err := h.service.GetFilmWithActorsAndCategoriesByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Film not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filmWithActors)
}
