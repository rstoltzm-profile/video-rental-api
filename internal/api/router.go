package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rstoltzm-profile/video-rental-api/internal/customer"
	"github.com/rstoltzm-profile/video-rental-api/internal/db"
	"github.com/rstoltzm-profile/video-rental-api/internal/film"
	"github.com/rstoltzm-profile/video-rental-api/internal/inventory"
	"github.com/rstoltzm-profile/video-rental-api/internal/rental"
	"github.com/rstoltzm-profile/video-rental-api/internal/store"
)

func NewRouter(pool *pgxpool.Pool, apiKey string) http.Handler {
	mux := http.NewServeMux()

	// health check
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/health/pool", healthHandlerWithPool(pool))

	// v1 routes
	v1 := http.NewServeMux()
	registerCustomerRoutes(v1, pool)
	registerRentalRoutes(v1, pool)
	registerInventoryRoutes(v1, pool)
	registerStoreRoutes(v1, pool)
	registerFilmRoutes(v1, pool)

	mux.Handle("/v1/", http.StripPrefix("/v1",
		requestSizeMiddleware(
			apiKeyMiddleware(apiKey,
				errorMiddleware(v1.ServeHTTP)))))
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

func registerCustomerRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := customer.NewRepository(pool)
	svc := customer.NewService(repo, repo, repo)
	handler := customer.NewHandler(svc)
	mux.HandleFunc("GET /customers", handler.GetCustomers)
	mux.HandleFunc("GET /customers/{id}", handler.GetCustomerByID)
	mux.HandleFunc("POST /customers", handler.CreateCustomer)
	mux.HandleFunc("DELETE /customers/{id}", handler.DeleteCustomerByID)
}

func registerRentalRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := rental.NewRepository(pool)
	svc := rental.NewService(repo, repo, repo)
	handler := rental.NewHandler(svc)
	mux.HandleFunc("GET /rentals", handler.GetRentals)
	mux.HandleFunc("POST /rentals", handler.CreateRental)
	mux.HandleFunc("POST /rentals/{id}/return", handler.ReturnRental)
}

func registerInventoryRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := inventory.NewRepository(pool)
	svc := inventory.NewService(repo, repo)
	handler := inventory.NewHandler(svc)
	mux.HandleFunc("GET /inventory", handler.GetInventory)
	mux.HandleFunc("GET /inventory/available", handler.GetInventoryAvailable)
}

func registerStoreRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := store.NewRepository(pool)
	svc := store.NewService(repo, repo)
	handler := store.NewHandler(svc)
	mux.HandleFunc("GET /stores/{id}/inventory/summary", handler.GetStoreInventorySummary)
}

func registerFilmRoutes(mux *http.ServeMux, pool *pgxpool.Pool) {
	repo := film.NewRepository(pool)
	svc := film.NewService(repo, repo)
	handler := film.NewHandler(svc)
	mux.HandleFunc("GET /films", handler.GetFilms)
	mux.HandleFunc("GET /films/{id}", handler.GetFilmByID)
	mux.HandleFunc("GET /films/search", handler.SearchFilm)
	mux.HandleFunc("GET /films/", handler.GetFilmWithActorsAndCategoriesByID)
}

func errorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
		log.Printf("Response completed for: %s %s", r.Method, r.URL.Path)
	}
}

func apiKeyMiddleware(validAPIKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for API key in header
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}

		if apiKey != validAPIKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		// API key is valid, continue
		next.ServeHTTP(w, r)
	}
}

func requestSizeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	const maxRequestSize = 1 << 20 // 1MB limit

	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)

		next.ServeHTTP(w, r)
	}
}

func healthHandlerWithPool(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check database health
		dbStatus := "ok"
		if err := db.HealthCheck(pool); err != nil {
			dbStatus = "unhealthy"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		resp := map[string]interface{}{
			"status":   "ok",
			"database": dbStatus,
			"connections": map[string]interface{}{
				"total":    pool.Stat().TotalConns(),
				"idle":     pool.Stat().IdleConns(),
				"acquired": pool.Stat().AcquiredConns(),
			},
		}

		json.NewEncoder(w).Encode(resp)
	}
}
