package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// configure server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      api.NewRouter(conn, cfg.APIKey),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// start the server, pass in conn to NewRouter so api can use it.
	log.Printf("Starting server on port %s...", cfg.Port)
	return server.ListenAndServe()
}
