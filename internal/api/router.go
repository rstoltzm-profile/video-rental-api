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

	// v1 routes
	v1 := http.NewServeMux()

	// customer routes
	repo := customer.NewRepository(conn)
	svc := customer.NewService(repo)
	handler := customer.NewHandler(svc)
	v1.HandleFunc("/customers", handler.GetCustomers)
	v1.HandleFunc("/customer/{id}", handler.GetCustomerByID)

	// mount v1 under /v1/
	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

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
