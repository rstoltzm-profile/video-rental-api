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
