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

type TVShowRepositoryIntegrationTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo TVShowRepository
}

func (suite *TVShowRepositoryIntegrationTestSuite) SetupSuite() {
	suite.db = SetupTestDB(suite.T())
	suite.repo = NewTVShowRepository(suite.db)
}

func (suite *TVShowRepositoryIntegrationTestSuite) TearDownTest() {
	CleanupTestData(suite.db)
}

func (suite *TVShowRepositoryIntegrationTestSuite) TestCreate() {
	tvShow := &models.TVShow{
		Title:        "Breaking Bad",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}

	err := suite.repo.Create(context.Background(), tvShow)

	suite.NoError(err)
	suite.NotZero(tvShow.ID)
	suite.Equal("Breaking Bad", tvShow.Title)
}

func (suite *TVShowRepositoryIntegrationTestSuite) TestGetByID() {
	tvShow := &models.TVShow{
		Title:        "Stranger Things",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), tvShow)

	result, err := suite.repo.GetByID(context.Background(), tvShow.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(tvShow.ID, result.ID)
	suite.Equal("Stranger Things", result.Title)
}

func (suite *TVShowRepositoryIntegrationTestSuite) TestGetByID_NotFound() {
	_, err := suite.repo.GetByID(context.Background(), 999)

	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *TVShowRepositoryIntegrationTestSuite) TestUpdate() {
	tvShow := &models.TVShow{
		Title:        "Original Title",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), tvShow)

	tvShow.Title = "Updated Title"
	locale := "en-US"
	tvShow.Locale = &locale
	err := suite.repo.Update(context.Background(), tvShow)

	suite.NoError(err)

	result, _ := suite.repo.GetByID(context.Background(), tvShow.ID)
	suite.Equal("Updated Title", result.Title)
	suite.Equal(&locale, result.Locale)
}

func (suite *TVShowRepositoryIntegrationTestSuite) TestDelete() {
	tvShow := &models.TVShow{
		Title:        "To Be Deleted",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), tvShow)

	err := suite.repo.Delete(context.Background(), tvShow.ID)
	suite.NoError(err)

	_, err = suite.repo.GetByID(context.Background(), tvShow.ID)
	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *TVShowRepositoryIntegrationTestSuite) TestList() {
	for i := 1; i <= 15; i++ {
		tvShow := &models.TVShow{
			Title:        "TV Show " + string(rune(i)),
			CreatedDate:  time.Now(),
			ModifiedDate: time.Now(),
		}
		suite.repo.Create(context.Background(), tvShow)
	}

	tvShows, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(15), total)
	suite.Len(tvShows, 10)
}

func (suite *TVShowRepositoryIntegrationTestSuite) TestList_EmptyDB() {
	tvShows, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(0), total)
	suite.Len(tvShows, 0)
}

func (suite *TVShowRepositoryIntegrationTestSuite) TestUniqueConstraint() {
	tvShow1 := &models.TVShow{
		Title:        "Unique Title",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	err := suite.repo.Create(context.Background(), tvShow1)
	suite.NoError(err)

	tvShow2 := &models.TVShow{
		Title:        "Unique Title",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	err = suite.repo.Create(context.Background(), tvShow2)

	suite.Error(err)
}

func TestTVShowRepositoryIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, new(TVShowRepositoryIntegrationTestSuite))
}
