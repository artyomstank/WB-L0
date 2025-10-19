package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Setup test environment
	os.Setenv("HTTP_HOST", "testhost")
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("HTTP_TIMEOUT", "10s")
	os.Setenv("POSTGRES_HOST", "testdb")
	os.Setenv("CACHE_STARTUP_SIZE", "100")

	cfg := LoadConfig()

	assert.Equal(t, "testhost", cfg.HTTPServer.Host)
	assert.Equal(t, 8080, cfg.HTTPServer.Port)
	assert.Equal(t, 10*time.Second, cfg.HTTPServer.Timeout)
	assert.Equal(t, "testdb", cfg.Postgres.Host)
	assert.Equal(t, 100, cfg.Cache.StartupSize)

	// Test default values
	os.Clearenv()
	cfg = LoadConfig()

	assert.Equal(t, "localhost", cfg.HTTPServer.Host) // default value
	assert.Equal(t, 8081, cfg.HTTPServer.Port)        // default value
}

func TestGetDBConnStr(t *testing.T) {
	cfg := Postgres{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "pass",
		Database: "testdb",
	}

	expected := "host=localhost port=5432 user=test password=pass dbname=testdb sslmode=disable"
	assert.Equal(t, expected, cfg.GetDBConnStr())
}
