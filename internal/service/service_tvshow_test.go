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

type TVShowServiceTestSuite struct {
	suite.Suite
	repo    *mocks.MockTVShowRepository
	service *TVShowService
}

func (suite *TVShowServiceTestSuite) SetupTest() {
	suite.repo = new(mocks.MockTVShowRepository)
	suite.service = NewTVShowService(suite.repo)
}

func (suite *TVShowServiceTestSuite) TestCreate() {
	tvShow := &models.TVShow{Title: "Test TV Show"}
	suite.repo.On("Create", mock.Anything, tvShow).Return(nil)

	err := suite.service.Create(context.Background(), tvShow)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestCreate_Error() {
	tvShow := &models.TVShow{Title: "Test TV Show"}
	suite.repo.On("Create", mock.Anything, tvShow).Return(errors.New("database error"))

	err := suite.service.Create(context.Background(), tvShow)
	suite.Error(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestGetByID() {
	tvShow := &models.TVShow{ID: 1, Title: "Test TV Show"}
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(tvShow, nil)

	result, err := suite.service.GetByID(context.Background(), 1)
	suite.NoError(err)
	suite.Equal(tvShow, result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestGetByID_Error() {
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("not found"))

	result, err := suite.service.GetByID(context.Background(), 1)
	suite.Error(err)
	suite.Nil(result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestUpdate() {
	tvShow := &models.TVShow{ID: 1, Title: "Updated TV Show"}
	suite.repo.On("Update", mock.Anything, tvShow).Return(nil)

	err := suite.service.Update(context.Background(), tvShow)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestUpdate_Error() {
	tvShow := &models.TVShow{ID: 1, Title: "Updated TV Show"}
	suite.repo.On("Update", mock.Anything, tvShow).Return(errors.New("database error"))

	err := suite.service.Update(context.Background(), tvShow)
	suite.Error(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestDelete() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := suite.service.Delete(context.Background(), 1)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestDelete_Error() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error"))

	err := suite.service.Delete(context.Background(), 1)
	suite.Error(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestList() {
	tvShows := []models.TVShow{
		{ID: 1, Title: "TV Show 1"},
		{ID: 2, Title: "TV Show 2"},
	}
	suite.repo.On("List", mock.Anything, 1, 10).Return(tvShows, int64(2), nil)

	result, total, err := suite.service.List(context.Background(), 1, 10)
	suite.NoError(err)
	suite.Equal(tvShows, result)
	suite.Equal(int64(2), total)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestList_DefaultPagination() {
	tvShows := []models.TVShow{
		{ID: 1, Title: "TV Show 1"},
	}
	// When page < 1, it should default to 1
	// When limit < 1, it should default to 10
	suite.repo.On("List", mock.Anything, 1, 10).Return(tvShows, int64(1), nil)

	result, total, err := suite.service.List(context.Background(), 0, 0)
	suite.NoError(err)
	suite.Equal(tvShows, result)
	suite.Equal(int64(1), total)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowServiceTestSuite) TestList_Error() {
	suite.repo.On("List", mock.Anything, 1, 10).Return(nil, int64(0), errors.New("database error"))

	result, total, err := suite.service.List(context.Background(), 1, 10)
	suite.Error(err)
	suite.Nil(result)
	suite.Equal(int64(0), total)
	suite.repo.AssertExpectations(suite.T())
}

func TestTVShowServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TVShowServiceTestSuite))
}
