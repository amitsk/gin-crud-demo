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

type ViewSummaryRepositoryIntegrationTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repo       ViewSummaryRepository
	movieRepo  MovieRepository
	seasonRepo SeasonRepository
	testMovie  *models.Movie
	testSeason *models.Season
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) SetupSuite() {
	suite.db = SetupTestDB(suite.T())
	suite.repo = NewViewSummaryRepository(suite.db)
	suite.movieRepo = NewMovieRepository(suite.db)
	suite.seasonRepo = NewSeasonRepository(suite.db)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) SetupTest() {
	// Create test movie
	suite.testMovie = &models.Movie{
		Title:        "Test Movie",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.movieRepo.Create(context.Background(), suite.testMovie)

	// Create test season
	suite.testSeason = &models.Season{
		Title:        "Test Season",
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.seasonRepo.Create(context.Background(), suite.testSeason)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TearDownTest() {
	CleanupTestData(suite.db)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestCreate() {
	viewSummary := &models.ViewSummary{
		Duration:     "WEEKLY",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 0, 7),
		HoursViewed:  1000,
		MovieID:      &suite.testMovie.ID,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}

	err := suite.repo.Create(context.Background(), viewSummary)

	suite.NoError(err)
	suite.NotZero(viewSummary.ID)
	suite.Equal("WEEKLY", viewSummary.Duration)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestCreate_WithSeason() {
	viewSummary := &models.ViewSummary{
		Duration:     "SEMI_ANNUALLY",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 6, 0),
		HoursViewed:  5000,
		SeasonID:     &suite.testSeason.ID,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}

	err := suite.repo.Create(context.Background(), viewSummary)

	suite.NoError(err)
	suite.NotZero(viewSummary.ID)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestGetByID() {
	viewSummary := &models.ViewSummary{
		Duration:     "WEEKLY",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 0, 7),
		HoursViewed:  1000,
		MovieID:      &suite.testMovie.ID,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), viewSummary)

	result, err := suite.repo.GetByID(context.Background(), viewSummary.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(viewSummary.ID, result.ID)
	suite.Equal("WEEKLY", result.Duration)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestGetByID_WithMovie() {
	viewSummary := &models.ViewSummary{
		Duration:     "WEEKLY",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 0, 7),
		HoursViewed:  1000,
		MovieID:      &suite.testMovie.ID,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), viewSummary)

	result, err := suite.repo.GetByID(context.Background(), viewSummary.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.NotNil(result.Movie)
	suite.Equal(suite.testMovie.ID, result.Movie.ID)
	suite.Equal("Test Movie", result.Movie.Title)
	suite.Nil(result.Season)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestGetByID_WithSeason() {
	viewSummary := &models.ViewSummary{
		Duration:     "WEEKLY",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 0, 7),
		HoursViewed:  1000,
		SeasonID:     &suite.testSeason.ID,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), viewSummary)

	result, err := suite.repo.GetByID(context.Background(), viewSummary.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.NotNil(result.Season)
	suite.Equal(suite.testSeason.ID, result.Season.ID)
	suite.Equal("Test Season", result.Season.Title)
	suite.Nil(result.Movie)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestGetByID_NotFound() {
	_, err := suite.repo.GetByID(context.Background(), 999)

	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestUpdate() {
	viewSummary := &models.ViewSummary{
		Duration:     "WEEKLY",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 0, 7),
		HoursViewed:  1000,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), viewSummary)

	viewSummary.HoursViewed = 2000
	viewRank := 5
	viewSummary.ViewRank = &viewRank
	err := suite.repo.Update(context.Background(), viewSummary)

	suite.NoError(err)

	result, _ := suite.repo.GetByID(context.Background(), viewSummary.ID)
	suite.Equal(2000, result.HoursViewed)
	suite.Equal(&viewRank, result.ViewRank)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestDelete() {
	viewSummary := &models.ViewSummary{
		Duration:     "WEEKLY",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 0, 7),
		HoursViewed:  1000,
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
	}
	suite.repo.Create(context.Background(), viewSummary)

	err := suite.repo.Delete(context.Background(), viewSummary.ID)
	suite.NoError(err)

	_, err = suite.repo.GetByID(context.Background(), viewSummary.ID)
	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestList() {
	for i := 1; i <= 15; i++ {
		viewSummary := &models.ViewSummary{
			Duration:     "WEEKLY",
			StartDate:    time.Now(),
			EndDate:      time.Now().AddDate(0, 0, 7),
			HoursViewed:  1000 + i,
			MovieID:      &suite.testMovie.ID,
			CreatedDate:  time.Now(),
			ModifiedDate: time.Now(),
		}
		suite.repo.Create(context.Background(), viewSummary)
	}

	summaries, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(15), total)
	suite.Len(summaries, 10)
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestList_WithRelations() {
	for i := 1; i <= 5; i++ {
		viewSummary := &models.ViewSummary{
			Duration:     "WEEKLY",
			StartDate:    time.Now(),
			EndDate:      time.Now().AddDate(0, 0, 7),
			HoursViewed:  1000 + i,
			MovieID:      &suite.testMovie.ID,
			CreatedDate:  time.Now(),
			ModifiedDate: time.Now(),
		}
		suite.repo.Create(context.Background(), viewSummary)
	}

	summaries, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(5), total)
	suite.Len(summaries, 5)

	// Verify preloading worked
	for _, summary := range summaries {
		suite.NotNil(summary.Movie)
		suite.Equal("Test Movie", summary.Movie.Title)
	}
}

func (suite *ViewSummaryRepositoryIntegrationTestSuite) TestList_EmptyDB() {
	summaries, total, err := suite.repo.List(context.Background(), 1, 10)

	suite.NoError(err)
	suite.Equal(int64(0), total)
	suite.Len(summaries, 0)
}

func TestViewSummaryRepositoryIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, new(ViewSummaryRepositoryIntegrationTestSuite))
}
