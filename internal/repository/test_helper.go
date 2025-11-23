//go:build integration

package repository

import (
	"testing"

	"github.com/amit/gin-crud-demo/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.Movie{},
		&models.TVShow{},
		&models.Season{},
		&models.ViewSummary{},
	)
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

// CleanupTestData truncates all tables
func CleanupTestData(db *gorm.DB) {
	db.Exec("DELETE FROM view_summary")
	db.Exec("DELETE FROM season")
	db.Exec("DELETE FROM tv_show")
	db.Exec("DELETE FROM movie")
}
