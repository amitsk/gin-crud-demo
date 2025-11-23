package service

import (
	"context"
	"errors"
	"testing"

	"github.com/amit/gin-crud-demo/internal/models"
	"github.com/amit/gin-crud-demo/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SeasonServiceTestSuite struct {
	suite.Suite
	repo    *mocks.MockSeasonRepository
	service *SeasonService
}

func (suite *SeasonServiceTestSuite) SetupTest() {
	suite.repo = new(mocks.MockSeasonRepository)
	suite.service = NewSeasonService(suite.repo)
}

func (suite *SeasonServiceTestSuite) TestCreate() {
	season := &models.Season{Title: "Test Season"}
	suite.repo.On("Create", mock.Anything, season).Return(nil)

	err := suite.service.Create(context.Background(), season)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestCreate_Error() {
	season := &models.Season{Title: "Test Season"}
	suite.repo.On("Create", mock.Anything, season).Return(errors.New("database error"))

	err := suite.service.Create(context.Background(), season)
	suite.Error(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestGetByID() {
	season := &models.Season{ID: 1, Title: "Test Season"}
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(season, nil)

	result, err := suite.service.GetByID(context.Background(), 1)
	suite.NoError(err)
	suite.Equal(season, result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestGetByID_Error() {
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("not found"))

	result, err := suite.service.GetByID(context.Background(), 1)
	suite.Error(err)
	suite.Nil(result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestUpdate() {
	season := &models.Season{ID: 1, Title: "Updated Season"}
	suite.repo.On("Update", mock.Anything, season).Return(nil)

	err := suite.service.Update(context.Background(), season)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestUpdate_Error() {
	season := &models.Season{ID: 1, Title: "Updated Season"}
	suite.repo.On("Update", mock.Anything, season).Return(errors.New("database error"))

	err := suite.service.Update(context.Background(), season)
	suite.Error(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestDelete() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := suite.service.Delete(context.Background(), 1)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestDelete_Error() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error"))

	err := suite.service.Delete(context.Background(), 1)
	suite.Error(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestList() {
	seasons := []models.Season{
		{ID: 1, Title: "Season 1"},
		{ID: 2, Title: "Season 2"},
	}
	suite.repo.On("List", mock.Anything, 1, 10).Return(seasons, int64(2), nil)

	result, total, err := suite.service.List(context.Background(), 1, 10)
	suite.NoError(err)
	suite.Equal(seasons, result)
	suite.Equal(int64(2), total)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestList_DefaultPagination() {
	seasons := []models.Season{
		{ID: 1, Title: "Season 1"},
	}
	// When page < 1, it should default to 1
	// When limit < 1, it should default to 10
	suite.repo.On("List", mock.Anything, 1, 10).Return(seasons, int64(1), nil)

	result, total, err := suite.service.List(context.Background(), 0, 0)
	suite.NoError(err)
	suite.Equal(seasons, result)
	suite.Equal(int64(1), total)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonServiceTestSuite) TestList_Error() {
	suite.repo.On("List", mock.Anything, 1, 10).Return(nil, int64(0), errors.New("database error"))

	result, total, err := suite.service.List(context.Background(), 1, 10)
	suite.Error(err)
	suite.Nil(result)
	suite.Equal(int64(0), total)
	suite.repo.AssertExpectations(suite.T())
}

func TestSeasonServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SeasonServiceTestSuite))
}
