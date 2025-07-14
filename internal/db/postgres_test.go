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

		if got == "" {
			t.Errorf("expected port env variable")
		}

	})
}
