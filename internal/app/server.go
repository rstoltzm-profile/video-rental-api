package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rstoltzm-profile/video-rental-api/internal/api"
	"github.com/rstoltzm-profile/video-rental-api/internal/config"
	"github.com/rstoltzm-profile/video-rental-api/internal/db"
)

func Run() error {
	// load configs
	cfg := config.LoadConfig()

	// Create connection pool with retry logic
	pool, err := connectWithRetry(cfg.DatabaseURL, 3)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	// Initial health check
	if err := db.HealthCheck(pool); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Run migrations
	if err := db.RunMigrations(pool); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Log current migration version
	version, appliedAt, err := db.GetCurrentMigration(pool)
	if err != nil {
		log.Printf("Could not fetch migration version: %v", err)
	} else {
		log.Printf("Database schema version: %s (applied at %s)", version, appliedAt)
	}

	// configure server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      api.NewRouter(pool, cfg.APIKey),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// start the server, pass in conn to NewRouter so api can use it.
	log.Printf("Starting server on port %s...", cfg.Port)
	return server.ListenAndServe()
}

func connectWithRetry(url string, maxRetries int) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < maxRetries; i++ {
		pool, err = db.ConnectPool(url)
		if err == nil {
			log.Printf("Connected to database on attempt %d", i+1)
			return pool, nil
		}

		log.Printf("Database connection attempt %d failed: %v", i+1, err)
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 2 * time.Second) // Exponential backoff
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}
