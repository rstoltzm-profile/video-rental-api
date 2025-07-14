package db

import (
	"testing"

	"github.com/rstoltzm-profile/video-rental-api/internal/config"
)

func TestEnv(t *testing.T) {
	t.Run("test db url", func(t *testing.T) {
		got := config.LoadConfig()

		if got.DatabaseURL == "" {
			t.Errorf("expected DatabaseURL env variable")
		}

	})

	t.Run("test db port", func(t *testing.T) {
		got := config.LoadConfig().Port
		want := "8080"

		if got != want {
			t.Errorf("expected port %v, got %v", want, got)
		}

	})
}
