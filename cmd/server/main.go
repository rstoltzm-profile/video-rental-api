package main

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

func main() {
	cfg := config.LoadConfig()
	conn, err := db.Connect(cfg.DatabaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)

	log.Printf("Starting server on port %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, api.NewRouter()))
}
