package db

import (
	"context"
	"testing"

	"github.com/rstoltzm-profile/video-rental-api/internal/config"
)

func TestConn(t *testing.T) {
	cfg := config.LoadConfig()
	conn, err := Connect(cfg.DatabaseURL)

	if cfg.DatabaseURL == "" {
		t.Skip("Skipping test: DATABASE_URL not set")
	}

	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		t.Fatalf("QueryRow failed: %v", err)
	}

	t.Logf("DB responded with: %s", greeting)
}
