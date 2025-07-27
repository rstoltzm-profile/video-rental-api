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

// GetInventory godoc
// @Summary      Get inventory
// @Description  Retrieve inventory items. Optionally filter by store_id query param.
// @Tags         inventory
// @Produce      json
// @Param        store_id  query     int     false  "Store ID to filter inventory"
// @Success      200       {array}   inventory.Inventory
// @Failure      500       {string}  string  "Internal Server Error"
// @Router       /inventory [get]
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

// GetInventoryAvailable godoc
// @Summary      Check inventory availability
// @Description  Check if a specific film is available in a given store
// @Tags         inventory
// @Produce      json
// @Param        store_id  query     int     true   "Store ID"
// @Param        film_id   query     int     true   "Film ID"
// @Success      200       {object}  inventory.InventoryAvailability
// @Failure      404       {object}  map[string]interface{}  "Not Found, available=false"
// @Failure      500       {string}  string  "Internal Server Error"
// @Router       /inventory/available [get]
func (h *Handler) GetInventoryAvailable(w http.ResponseWriter, r *http.Request) {
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
