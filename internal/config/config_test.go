package config_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/tktanisha/booking_system/internal/config"
)

func TestLoadEnv(t *testing.T) {
	oldEnv := os.Environ()
	defer func() {
		for _, e := range oldEnv {
			parts := splitEnv(e)
			os.Setenv(parts[0], parts[1])
		}
	}()

	os.Unsetenv("DATABASE_URL")

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	config.LoadEnv()

	if logBuf.Len() == 0 {
		t.Errorf("expected log output when .env is missing, got empty")
	}
}

func TestGetDBURL_Success(t *testing.T) {
	expected := "postgres://user:pass@localhost:5432/dbname"
	os.Setenv("DATABASE_URL", expected)
	defer os.Unsetenv("DATABASE_URL")

	got := config.GetDBURL()
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestGetDBURL(t *testing.T) {
	// Save original env value
	original := os.Getenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", original)

	t.Run("Returns DB URL when set", func(t *testing.T) {
		expected := "postgres://user:pass@localhost:5432/dbname"
		os.Setenv("DATABASE_URL", expected)

		got := config.GetDBURL()
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("Returns empty string when not set", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")

		got := config.GetDBURL()
		if got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})
}

func splitEnv(env string) [2]string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return [2]string{env[:i], env[i+1:]}
		}
	}
	return [2]string{env, ""}
}
