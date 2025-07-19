package inventory

import (
	"database/sql"
	"encoding/json"
	"errors"
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

func (h *Handler) GetInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	storeIDStr := r.URL.Query().Get("store_id")

	var inventory []Inventory
	var err error
	storeID := -1

	if storeIDStr != "" {
		storeID, _ = strconv.Atoi(storeIDStr)
		inventory, err = h.service.GetInventoryByStore(r.Context(), storeID)
	} else {
		inventory, err = h.service.GetInventory(r.Context())
	}

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(inventory)
}

func (h *Handler) GetInventoryAvailable(w http.ResponseWriter, r *http.Request) {
	// GET /inventory/available?film_id=1&store_id=2
	w.Header().Set("Content-Type", "application/json")

	var err error
	var inventoryAvailability InventoryAvailability
	storeIDStr := r.URL.Query().Get("store_id")
	filmIDStr := r.URL.Query().Get("film_id")
	storeID := -1
	filmID := -1

	if storeIDStr != "" {
		storeID, err = strconv.Atoi(storeIDStr)
	}

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	if filmIDStr != "" {
		filmID, err = strconv.Atoi(filmIDStr)
	}

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	if storeID != -1 && filmID != -1 {
		inventoryAvailability, err = h.service.GetInventoryAvailable(r.Context(), storeID, filmID)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]any{
				"store_id":  storeID,
				"film_id":   filmID,
				"available": false,
			})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(inventoryAvailability)

}
