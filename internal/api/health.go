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

func landingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Video Rental API</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					color: #333;
					text-align: center;
					padding: 50px;
				}
				h1 {
					color: #2c3e50;
				}
				ul {
					list-style: none;
					padding: 0;
				}
				li {
					margin: 15px 0;
				}
				a {
					text-decoration: none;
					color: #3498db;
					font-size: 18px;
					font-weight: bold;
				}
				a:hover {
					color: #1abc9c;
				}
				.container {
					background: white;
					padding: 30px;
					border-radius: 8px;
					box-shadow: 0 4px 8px rgba(0,0,0,0.1);
					display: inline-block;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Welcome to the Video Rental API</h1>
				<p>Select an option below:</p>
				<ul>
					<li><a href="/swagger/index.html">📜 Swagger API Docs</a></li>
					<li><a href="/health">✅ Health Check</a></li>
				</ul>
			</div>
		</body>
		</html>
	`))
}
