package rental

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

func (h *Handler) GetRentals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parameters
	late := r.URL.Query().Get("late")
	customerIDStr := r.URL.Query().Get("customer_id")

	var rentals []Rental
	var err error
	customerID := -1

	if customerIDStr != "" {
		customerID, err = strconv.Atoi(customerIDStr)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}
	}

	if customerID != -1 && late == "true" {
		rentals, err = h.service.GetLateRentalsByCustomerID(r.Context(), customerID)
	} else if customerID != -1 {
		rentals, err = h.service.GetRentalsByCustomerID(r.Context(), customerID)
	} else if late == "true" {
		rentals, err = h.service.GetRentals(r.Context())
	} else {
		rentals, err = h.service.GetLateRentals(r.Context())
	}

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rentals)
}

func (h *Handler) CreateRental(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateRentalRequest

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
