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

type MovieServiceTestSuite struct {
	suite.Suite
	repo    *mocks.MockMovieRepository
	service *MovieService
}

func (suite *MovieServiceTestSuite) SetupTest() {
	suite.repo = new(mocks.MockMovieRepository)
	suite.service = NewMovieService(suite.repo)
}

func (suite *MovieServiceTestSuite) TestCreate() {
	movie := &models.Movie{Title: "Test Movie"}
	suite.repo.On("Create", mock.Anything, movie).Return(nil)

	err := suite.service.Create(context.Background(), movie)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestGetByID() {
	movie := &models.Movie{ID: 1, Title: "Test Movie"}
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(movie, nil)

	result, err := suite.service.GetByID(context.Background(), 1)
	suite.NoError(err)
	suite.Equal(movie, result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestGetByID_Error() {
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("not found"))

	result, err := suite.service.GetByID(context.Background(), 1)
	suite.Error(err)
	suite.Nil(result)
	suite.repo.AssertExpectations(suite.T())
}

func TestMovieServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MovieServiceTestSuite))
}
