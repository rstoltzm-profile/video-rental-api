package rental

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetRentals godoc
// @Summary      List rentals
// @Description  Get all rentals or only late rentals based on query param
// @Tags         rentals
// @Accept       json
// @Produce      json
// @Param        late  query     bool  false  "Filter late rentals (true)"
// @Success      200   {array}   rental.Rental
// @Failure      500   {string}  string  "Failed to fetch rentals"
// @Security     ApiKeyAuth
// @Router       /v1/rentals [get]
func (h *Handler) GetRentals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parameters
	late := r.URL.Query().Get("late")

	var rentals []Rental
	var err error

	if late == "true" {
		rentals, err = h.service.GetLateRentals(r.Context())
	} else {
		rentals, err = h.service.GetRentals(r.Context())
	}

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rentals)
}

// CreateRental godoc
// @Summary      Create a rental
// @Description  Create a new rental record
// @Tags         rentals
// @Accept       json
// @Produce      json
// @Param        rental  body      rental.CreateRentalRequest  true  "Rental request"
// @Success      201     {object}  map[string]int  "Created rental ID"
// @Failure      400     {string}  string  "Invalid input"
// @Failure      500     {string}  string  "Failed to create rental"
// @Security     ApiKeyAuth
// @Router       /v1/rentals [post]
func (h *Handler) CreateRental(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateRentalRequest

	// Validate the request
	if err := validate.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	rental, err := h.service.CreateRental(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create rental", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/rentals/%d", rental))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": rental})
}

// ReturnRental godoc
// @Summary      Return rental
// @Description  Mark a rental as returned by ID
// @Tags         rentals
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Rental ID"
// @Success      204  {string}  string  "Rental returned successfully"
// @Failure      400  {string}  string  "Invalid rental ID"
// @Failure      500  {string}  string  "Failed to return rental"
// @Security     ApiKeyAuth
// @Router       /v1/rentals/{id}/return [put]
func (h *Handler) ReturnRental(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}
	err = h.service.ReturnRentalByID(r.Context(), id)

	if err != nil {
		http.Error(w, "Failed to return rental", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204
}
