package api

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rstoltzm-profile/video-rental-api/internal/customer"
)

func NewRouter(conn *pgx.Conn) http.Handler {
	mux := http.NewServeMux()

	// health check
	mux.HandleFunc("/health", healthHandler)

	// customer routes
	repo := customer.NewRepository(conn)
	svc := customer.NewService(repo)
	handler := customer.NewHandler(svc)
	mux.HandleFunc("/customers", handler.GetCustomers)
	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"status": "ok", // more common than "Health": "GOOD"
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode health response", http.StatusInternalServerError)
	}
}
