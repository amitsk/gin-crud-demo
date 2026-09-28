package handler

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amit/gin-crud-demo/internal/api/dto"
	"github.com/amit/gin-crud-demo/internal/models"
	"github.com/amit/gin-crud-demo/internal/service"
	"github.com/amit/gin-crud-demo/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ViewSummaryHandlerTestSuite struct {
	suite.Suite
	repo    *mocks.MockViewSummaryRepository
	service *service.ViewSummaryService
	handler *ViewSummaryHandler
	router  *gin.Engine
}

func (suite *ViewSummaryHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.repo = new(mocks.MockViewSummaryRepository)
	suite.service = service.NewViewSummaryService(suite.repo)
	suite.handler = NewViewSummaryHandler(suite.service)
	suite.router = gin.New()
	suite.router.POST("/view-summaries", suite.handler.Create)
	suite.router.GET("/view-summaries/:id", suite.handler.GetByID)
	suite.router.PUT("/view-summaries/:id", suite.handler.Update)
	suite.router.DELETE("/view-summaries/:id", suite.handler.Delete)
	suite.router.GET("/view-summaries", suite.handler.List)
}

func (suite *ViewSummaryHandlerTestSuite) TestCreate() {
	req := CreateViewSummaryRequest{
		Duration:    "WEEKLY",
		StartDate:   "2024-01-01",
		EndDate:     "2024-01-07",
		HoursViewed: 1000,
	}
	body, _ := json.Marshal(req)

	suite.repo.On("Create", mock.Anything, mock.AnythingOfType("*models.ViewSummary")).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/view-summaries", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusCreated, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryHandlerTestSuite) TestCreate_InvalidJSON() {
	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/view-summaries", bytes.NewBuffer([]byte("invalid json")))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *ViewSummaryHandlerTestSuite) TestCreate_MissingRequiredField() {
	req := CreateViewSummaryRequest{} // Missing required fields
	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/view-summaries", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *ViewSummaryHandlerTestSuite) TestCreate_InvalidDateFormat() {
	req := CreateViewSummaryRequest{
		Duration:    "WEEKLY",
		StartDate:   "invalid-date",
		EndDate:     "2024-01-07",
		HoursViewed: 1000,
	}
	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/view-summaries", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *ViewSummaryHandlerTestSuite) TestGetByID() {
	viewSummary := &models.ViewSummary{
		ID:          1,
		Duration:    "WEEKLY",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		HoursViewed: 1000,
	}
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(viewSummary, nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/view-summaries/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	var resp models.ViewSummary
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal(viewSummary.Duration, resp.Duration)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryHandlerTestSuite) TestGetByID_InvalidID() {
	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/view-summaries/invalid", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *ViewSummaryHandlerTestSuite) TestGetByID_NotFound() {
	suite.repo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/view-summaries/999", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryHandlerTestSuite) TestUpdate() {
	existingViewSummary := &models.ViewSummary{
		ID:          1,
		Duration:    "WEEKLY",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 7),
		HoursViewed: 1000,
	}
	req := CreateViewSummaryRequest{
		Duration:    "SEMI_ANNUALLY",
		StartDate:   "2024-01-01",
		EndDate:     "2024-06-30",
		HoursViewed: 5000,
	}
	body, _ := json.Marshal(req)

	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(existingViewSummary, nil)
	suite.repo.On("Update", mock.Anything, mock.AnythingOfType("*models.ViewSummary")).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("PUT", "/view-summaries/1", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryHandlerTestSuite) TestUpdate_NotFound() {
	req := CreateViewSummaryRequest{
		Duration:    "WEEKLY",
		StartDate:   "2024-01-01",
		EndDate:     "2024-01-07",
		HoursViewed: 1000,
	}
	body, _ := json.Marshal(req)

	suite.repo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("PUT", "/view-summaries/999", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryHandlerTestSuite) TestDelete() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("DELETE", "/view-summaries/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryHandlerTestSuite) TestDelete_Error() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("DELETE", "/view-summaries/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *ViewSummaryHandlerTestSuite) TestList() {
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

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/view-summaries?page=1&limit=10", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	var resp dto.ListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal(int64(2), resp.Total)
	suite.repo.AssertExpectations(suite.T())
}

func TestViewSummaryHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ViewSummaryHandlerTestSuite))
}
