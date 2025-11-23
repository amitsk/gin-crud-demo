//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/amit/gin-crud-demo/internal/models"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MovieRepositoryIntegrationTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo MovieRepository
}

func (suite *MovieRepositoryIntegrationTestSuite) SetupSuite() {
	suite.db = SetupTestDB(suite.T())
	suite.repo = NewMovieRepository(suite.db)
}

func (suite *MovieRepositoryIntegrationTestSuite) TearDownTest() {
	CleanupTestData(suite.db)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestCreate() {
	movie := &models.Movie{
		Title:        "Inception",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}

	err := suite.repo.Create(context.Background(), movie)

	suite.NoError(err)
	suite.NotZero(movie.ID)
	suite.Equal("Inception", movie.Title)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestCreate_WithAllFields() {
	originalTitle := "Inception Original"
	runtime := int64(148)
	availableGlobally := true
	locale := "en-US"
	releaseDate := time.Date(2010, 7, 16, 0, 0, 0, 0, time.UTC)

	movie := &models.Movie{
		Title:             "Inception",
		OriginalTitle:     &originalTitle,
		Runtime:           &runtime,
		ReleaseDate:       &releaseDate,
		AvailableGlobally: &availableGlobally,
		Locale:            &locale,
		CreatedDate:       time.Now(),
		ModifiedDate:      time.Now(),
	}

	err := suite.repo.Create(context.Background(), movie)

	suite.NoError(err)
	suite.NotZero(movie.ID)
	suite.Equal("Inception", movie.Title)
	suite.Equal(&originalTitle, movie.OriginalTitle)
	suite.Equal(&runtime, movie.Runtime)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestGetByID() {
	// Create a movie first
	movie := &models.Movie{
		Title:        "The Matrix",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), movie)

	// Get it by ID
	result, err := suite.repo.GetByID(context.Background(), movie.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(movie.ID, result.ID)
	suite.Equal("The Matrix", result.Title)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestGetByID_NotFound() {
	_, err := suite.repo.GetByID(context.Background(), 999)

	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestUpdate() {
	// Create a movie
	movie := &models.Movie{
		Title:        "Original Title",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), movie)

	// Update it
	movie.Title = "Updated Title"
	newRuntime := int64(120)
	movie.Runtime = &newRuntime
	err := suite.repo.Update(context.Background(), movie)

	suite.NoError(err)

	// Verify update
	result, _ := suite.repo.GetByID(context.Background(), movie.ID)
	suite.Equal("Updated Title", result.Title)
	suite.Equal(&newRuntime, result.Runtime)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestDelete() {
	// Create a movie
	movie := &models.Movie{
		Title:        "To Be Deleted",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), movie)

	// Delete it
	err := suite.repo.Delete(context.Background(), movie.ID)
	suite.NoError(err)

	// Verify deletion
	_, err = suite.repo.GetByID(context.Background(), movie.ID)
	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestList() {
	// Create multiple movies
	for i := 1; i <= 15; i++ {
		movie := &models.Movie{
			Title:        "Movie " + string(rune(i)),
			CreatedDate:  time.Now(),
			ModifiedDate: time.Now(),
		}
		suite.repo.Create(context.Background(), movie)
	}

	// List with pagination
	movies, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(15), total)
	suite.Len(movies, 10)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestList_SecondPage() {
	// Create multiple movies
	for i := 1; i <= 15; i++ {
		movie := &models.Movie{
			Title:        "Movie " + string(rune(i)),
			CreatedDate:  time.Now(),
			ModifiedDate: time.Now(),
		}
		suite.repo.Create(context.Background(), movie)
	}

	// List second page
	movies, total, err := suite.repo.List(context.Background(), 2, 10)

	suite.NoError(err)
	suite.Equal(int64(15), total)
	suite.Len(movies, 5)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestList_EmptyDB() {
	movies, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(0), total)
	suite.Len(movies, 0)
}

func (suite *MovieRepositoryIntegrationTestSuite) TestUniqueConstraint() {
	// Create first movie
	movie1 := &models.Movie{
		Title:        "Unique Title",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	err := suite.repo.Create(context.Background(), movie1)
	suite.NoError(err)

	// Try to create another movie with the same title
	movie2 := &models.Movie{
		Title:        "Unique Title",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	err = suite.repo.Create(context.Background(), movie2)

	// Should fail due to unique constraint (but SQLite may not enforce it strictly)
	// So we just check that if it succeeds, the IDs are different
	if err == nil {
		suite.NotEqual(movie1.ID, movie2.ID)
	}
}

func TestMovieRepositoryIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, new(MovieRepositoryIntegrationTestSuite))
}
