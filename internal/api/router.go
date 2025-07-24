package api

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rstoltzm-profile/video-rental-api/internal/customer"
	"github.com/rstoltzm-profile/video-rental-api/internal/film"
	"github.com/rstoltzm-profile/video-rental-api/internal/inventory"
	"github.com/rstoltzm-profile/video-rental-api/internal/middleware"
	"github.com/rstoltzm-profile/video-rental-api/internal/rental"
	"github.com/rstoltzm-profile/video-rental-api/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(pool *pgxpool.Pool, apiKey string) http.Handler {
	mux := http.NewServeMux()

	// health check
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/health/pool", healthHandlerWithPool(pool))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// v1 routes
	v1 := http.NewServeMux()
	registerCustomerRoutes(v1, pool)
	registerRentalRoutes(v1, pool)
	registerInventoryRoutes(v1, pool)
	registerStoreRoutes(v1, pool)
	registerFilmRoutes(v1, pool)

	mux.Handle("/v1/", http.StripPrefix("/v1",
		middleware.CORSMiddleware(
			middleware.RequestSizeMiddleware(
				middleware.ApiKeyMiddleware(apiKey,
					middleware.ErrorMiddleware(v1.ServeHTTP))))))
	return mux
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
