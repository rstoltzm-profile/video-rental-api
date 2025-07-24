package api

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rstoltzm-profile/video-rental-api/internal/db"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"status": "ok", // more common than "Health": "GOOD"
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode health response", http.StatusInternalServerError)
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
