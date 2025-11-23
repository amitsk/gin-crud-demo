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

type SeasonRepositoryIntegrationTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repo       SeasonRepository
	tvShowRepo TVShowRepository
	testTVShow *models.TVShow
}

func (suite *SeasonRepositoryIntegrationTestSuite) SetupSuite() {
	suite.db = SetupTestDB(suite.T())
	suite.repo = NewSeasonRepository(suite.db)
	suite.tvShowRepo = NewTVShowRepository(suite.db)
}

func (suite *SeasonRepositoryIntegrationTestSuite) SetupTest() {
	// Create a test TV show for foreign key relationships
	suite.testTVShow = &models.TVShow{
		Title:        "Test TV Show",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.tvShowRepo.Create(context.Background(), suite.testTVShow)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TearDownTest() {
	CleanupTestData(suite.db)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestCreate() {
	season := &models.Season{
		Title:        "Season 1",
		TVShowID:     &suite.testTVShow.ID,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}

	err := suite.repo.Create(context.Background(), season)

	suite.NoError(err)
	suite.NotZero(season.ID)
	suite.Equal("Season 1", season.Title)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestCreate_WithoutTVShow() {
	season := &models.Season{
		Title:        "Standalone Season",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}

	err := suite.repo.Create(context.Background(), season)

	suite.NoError(err)
	suite.NotZero(season.ID)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestGetByID() {
	season := &models.Season{
		Title:        "Season 1",
		TVShowID:     &suite.testTVShow.ID,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), season)

	result, err := suite.repo.GetByID(context.Background(), season.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(season.ID, result.ID)
	suite.Equal("Season 1", result.Title)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestGetByID_WithTVShow() {
	season := &models.Season{
		Title:        "Season 1",
		TVShowID:     &suite.testTVShow.ID,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), season)

	result, err := suite.repo.GetByID(context.Background(), season.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.NotNil(result.TVShow)
	suite.Equal(suite.testTVShow.ID, result.TVShow.ID)
	suite.Equal("Test TV Show", result.TVShow.Title)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestGetByID_NotFound() {
	_, err := suite.repo.GetByID(context.Background(), 999)

	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestUpdate() {
	season := &models.Season{
		Title:        "Original Title",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), season)

	season.Title = "Updated Title"
	seasonNumber := 2
	season.SeasonNumber = &seasonNumber
	err := suite.repo.Update(context.Background(), season)

	suite.NoError(err)

	result, _ := suite.repo.GetByID(context.Background(), season.ID)
	suite.Equal("Updated Title", result.Title)
	suite.Equal(&seasonNumber, result.SeasonNumber)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestDelete() {
	season := &models.Season{
		Title:        "To Be Deleted",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), season)

	err := suite.repo.Delete(context.Background(), season.ID)
	suite.NoError(err)

	_, err = suite.repo.GetByID(context.Background(), season.ID)
	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestList() {
	for i := 1; i <= 15; i++ {
		season := &models.Season{
			Title:        "Season " + string(rune(i)),
			TVShowID:     &suite.testTVShow.ID,
			CreatedDate:  time.Now(),
			ModifiedDate: time.Now(),
		}
		suite.repo.Create(context.Background(), season)
	}

	seasons, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(15), total)
	suite.Len(seasons, 10)
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestList_WithTVShow() {
	for i := 1; i <= 5; i++ {
		season := &models.Season{
			Title:        "Season " + string(rune(i)),
			TVShowID:     &suite.testTVShow.ID,
			CreatedDate:  time.Now(),
			ModifiedDate: time.Now(),
		}
		suite.repo.Create(context.Background(), season)
	}

	seasons, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(5), total)
	suite.Len(seasons, 5)

	// Verify preloading worked
	for _, season := range seasons {
		suite.NotNil(season.TVShow)
		suite.Equal("Test TV Show", season.TVShow.Title)
	}
}

func (suite *SeasonRepositoryIntegrationTestSuite) TestList_EmptyDB() {
	seasons, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(0), total)
	suite.Len(seasons, 0)
}

func TestSeasonRepositoryIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, new(SeasonRepositoryIntegrationTestSuite))
}
