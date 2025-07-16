package inventory

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
