package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/amit/gin-crud-demo/internal/models"
	"github.com/amit/gin-crud-demo/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ViewSummaryServiceTestSuite struct {
	suite.Suite
	repo    *mocks.MockViewSummaryRepository
	service *ViewSummaryService
}

func (suite *ViewSummaryServiceTestSuite) SetupTest() {
	suite.repo = new(mocks.MockViewSummaryRepository)
	suite.service = NewViewSummaryService(suite.repo)
}

func (suite *ViewSummaryServiceTestSuite) TestCreate() {
	viewSummary := &models.ViewSummary{
		Duration:    "WEEKLY",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		HoursViewed: 1000,
	}
	suite.repo.On("Create", mock.Anything, viewSummary).Return(nil)

	err := suite.service.Create(context.Background(), viewSummary)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestCreate_Error() {
	viewSummary := &models.ViewSummary{
		Duration:    "WEEKLY",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		HoursViewed: 1000,
	}
	suite.repo.On("Create", mock.Anything, viewSummary).Return(errors.New("database error"))

	err := suite.service.Create(context.Background(), viewSummary)
	suite.Error(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestGetByID() {
	viewSummary := &models.ViewSummary{
		ID:          1,
		Duration:    "WEEKLY",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		HoursViewed: 1000,
	}
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(viewSummary, nil)

	result, err := suite.service.GetByID(context.Background(), 1)
	suite.NoError(err)
	suite.Equal(viewSummary, result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestGetByID_Error() {
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("not found"))

	result, err := suite.service.GetByID(context.Background(), 1)
	suite.Error(err)
	suite.Nil(result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestUpdate() {
	viewSummary := &models.ViewSummary{
		ID:          1,
		Duration:    "WEEKLY",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		HoursViewed: 2000,
	}
	suite.repo.On("Update", mock.Anything, viewSummary).Return(nil)

	err := suite.service.Update(context.Background(), viewSummary)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestUpdate_Error() {
	viewSummary := &models.ViewSummary{
		ID:          1,
		Duration:    "WEEKLY",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		HoursViewed: 2000,
	}
	suite.repo.On("Update", mock.Anything, viewSummary).Return(errors.New("database error"))

	err := suite.service.Update(context.Background(), viewSummary)
	suite.Error(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestDelete() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := suite.service.Delete(context.Background(), 1)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestDelete_Error() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error"))

	err := suite.service.Delete(context.Background(), 1)
	suite.Error(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestList() {
	viewSummaries := []models.ViewSummary{
		{
			ID:          1,
			Duration:    "WEEKLY",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 0, 7),
			HoursViewed: 1000,
		},
		{
			ID:          2,
			Duration:    "SEMI_ANNUALLY",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 6, 0),
			HoursViewed: 5000,
		},
	}
	suite.repo.On("List", mock.Anything, 1, 10).Return(viewSummaries, int64(2), nil)

	result, total, err := suite.service.List(context.Background(), 1, 10)
	suite.NoError(err)
	suite.Equal(viewSummaries, result)
	suite.Equal(int64(2), total)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestList_DefaultPagination() {
	viewSummaries := []models.ViewSummary{
		{
			ID:          1,
			Duration:    "WEEKLY",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 0, 7),
			HoursViewed: 1000,
		},
	}
	// When page < 1, it should default to 1
	// When limit < 1, it should default to 10
	suite.repo.On("List", mock.Anything, 1, 10).Return(viewSummaries, int64(1), nil)

	result, total, err := suite.service.List(context.Background(), 0, 0)
	suite.NoError(err)
	suite.Equal(viewSummaries, result)
	suite.Equal(int64(1), total)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryServiceTestSuite) TestList_Error() {
	suite.repo.On("List", mock.Anything, 1, 10).Return(nil, int64(0), errors.New("database error"))

	result, total, err := suite.service.List(context.Background(), 1, 10)
	suite.Error(err)
	suite.Nil(result)
	suite.Equal(int64(0), total)
	suite.repo.AssertExpectations(suite.T())
}

func TestViewSummaryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ViewSummaryServiceTestSuite))
}
