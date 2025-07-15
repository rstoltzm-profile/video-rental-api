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
	registerCustomerRoutes(v1, conn)

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

func registerCustomerRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	repo := customer.NewRepository(conn)
	svc := customer.NewService(repo, repo, repo)
	handler := customer.NewHandler(svc)
	mux.HandleFunc("GET /customers", handler.GetCustomers)
	mux.HandleFunc("GET /customers/{id}", handler.GetCustomerByID)
	mux.HandleFunc("POST /customers", handler.CreateCustomer)
	mux.HandleFunc("DELETE /customers/{id}", handler.DeleteCustomerByID)
}
