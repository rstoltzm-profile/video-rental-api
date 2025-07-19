package api

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rstoltzm-profile/video-rental-api/internal/customer"
	"github.com/rstoltzm-profile/video-rental-api/internal/film"
	"github.com/rstoltzm-profile/video-rental-api/internal/inventory"
	"github.com/rstoltzm-profile/video-rental-api/internal/rental"
	"github.com/rstoltzm-profile/video-rental-api/internal/store"
)

func NewRouter(conn *pgx.Conn) http.Handler {
	mux := http.NewServeMux()

	// health check
	mux.HandleFunc("/health", healthHandler)

	// v1 routes
	v1 := http.NewServeMux()
	registerCustomerRoutes(v1, conn)
	registerRentalRoutes(v1, conn)
	registerInventoryRoutes(v1, conn)
	registerStoreRoutes(v1, conn)
	registerFilmRoutes(v1, conn)

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

func registerRentalRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	repo := rental.NewRepository(conn)
	svc := rental.NewService(repo, repo)
	handler := rental.NewHandler(svc)
	mux.HandleFunc("GET /rentals", handler.GetRentals)
}

func registerInventoryRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	repo := inventory.NewRepository(conn)
	svc := inventory.NewService(repo, repo)
	handler := inventory.NewHandler(svc)
	mux.HandleFunc("GET /inventory", handler.GetInventory)
	mux.HandleFunc("GET /inventory/available", handler.GetInventoryAvailable)
}

func registerStoreRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	repo := store.NewRepository(conn)
	svc := store.NewService(repo, repo)
	handler := store.NewHandler(svc)
	mux.HandleFunc("GET /stores/{id}/inventory/summary", handler.GetStoreInventorySummary)
}

func registerFilmRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	repo := film.NewRepository(conn)
	svc := film.NewService(repo, repo)
	handler := film.NewHandler(svc)
	mux.HandleFunc("GET /films", handler.GetFilms)
	mux.HandleFunc("GET /films/{id}", handler.GetFilmByID)
	mux.HandleFunc("GET /films/search", handler.SearchFilm)
	mux.HandleFunc("GET /films/", handler.GetFilmWithActorsAndCategoriesByID)
}
