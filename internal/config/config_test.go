package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Change CWD to project root so it can find config/
	// Assuming test is running in internal/config
	wd, _ := os.Getwd()
	if err := os.Chdir("../../"); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	// Test case 1: Load test config via APP_ENVIRONMENT
	os.Setenv("APP_ENVIRONMENT", "test")
	defer os.Unsetenv("APP_ENVIRONMENT")

	// Test case 2: Override a value via env var (requires APP_ prefix)
	os.Setenv("APP_DATABASE_USER", "env_user")
	defer os.Unsetenv("APP_DATABASE_USER")

	// Test case 3: Override another value via env var (requires APP_ prefix)
	os.Setenv("APP_DATABASE_DBNAME", "netflix_movies_test_db")
	defer os.Unsetenv("APP_DATABASE_DBNAME")

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check if values are loaded from config.test.yaml (DBName should be netflix_movies_test)
	assert.Equal(t, "netflix_movies_test_db", cfg.Database.DBName)

	// Check if env var override works
	assert.Equal(t, "env_user", cfg.Database.User)

	// Check if Environment is set correctly
	assert.Equal(t, "test", cfg.Environment)
}
