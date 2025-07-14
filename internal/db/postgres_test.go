package db

import (
	"testing"

	"github.com/rstoltzm-profile/video-rental-api/internal/config"
)

func TestEnv(t *testing.T) {
	got := config.LoadConfig()

	if got.DatabaseURL == "" {
		t.Errorf("expected DatabaseURL env variable")
	}

}
