package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_WithEnvVars(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test_user:pass@localhost/test_db")
	os.Setenv("PORT", "9999")
	defer os.Clearenv() // optional cleanup

	cfg := LoadConfig()

	assert.Equal(t, "postgres://test_user:pass@localhost/test_db", cfg.DatabaseURL)
	assert.Equal(t, "9999", cfg.Port)
}

func TestLoadConfig_DefaultPort(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://another:pass@localhost/other_db")
	os.Unsetenv("PORT")
	defer os.Clearenv()

	cfg := LoadConfig()

	assert.Equal(t, "postgres://another:pass@localhost/other_db", cfg.DatabaseURL)
	assert.Equal(t, "8080", cfg.Port) // default fallback
}
