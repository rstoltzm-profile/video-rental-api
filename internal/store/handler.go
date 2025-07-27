package store

import (
	"encoding/json"
	"fmt"
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

// GetStoreInventorySummary godoc
// @Summary      Get store inventory summary
// @Description  Returns a summary count of inventory for a given store ID
// @Tags         stores
// @Produce      json
// @Param        id   path      int  true  "Store ID"
// @Success      200  {array}   store.StoreInventorySummary
// @Failure      400  {string}  string "Invalid store ID"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /stores/{id}/inventory/summary [get]
func (h *Handler) GetStoreInventorySummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/stores/")
	parts := strings.Split(path, "/")

	storeIDStr := parts[0]

	var inventory []StoreInventorySummary
	var err error
	storeID := -1

	if storeIDStr != "" {
		storeID, _ = strconv.Atoi(storeIDStr)
		inventory, err = h.service.GetStoreInventorySummary(r.Context(), storeID)
	}

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(inventory)
}
