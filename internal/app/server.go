package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rstoltzm-profile/video-rental-api/internal/api"
	"github.com/rstoltzm-profile/video-rental-api/internal/config"
	"github.com/rstoltzm-profile/video-rental-api/internal/db"
)

func Run() error {
	// load configs
	cfg := config.LoadConfig()

	// attempt to connect to db
	conn, err := db.Connect(cfg.DatabaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// keep connection alive until Run() returns (server shutdown)
	defer conn.Close(context.Background())

	// start the server, pass in conn to NewRouter so api can use it.
	log.Printf("Starting server on port %s...", cfg.Port)
	return http.ListenAndServe(":"+cfg.Port, api.NewRouter(conn, cfg.APIKey))
}
