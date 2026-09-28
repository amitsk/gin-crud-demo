package handler

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amit/gin-crud-demo/internal/api/dto"
	"github.com/amit/gin-crud-demo/internal/models"
	"github.com/amit/gin-crud-demo/internal/service"
	"github.com/amit/gin-crud-demo/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SeasonHandlerTestSuite struct {
	suite.Suite
	repo    *mocks.MockSeasonRepository
	service *service.SeasonService
	handler *SeasonHandler
	router  *gin.Engine
}

func (suite *SeasonHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.repo = new(mocks.MockSeasonRepository)
	suite.service = service.NewSeasonService(suite.repo)
	suite.handler = NewSeasonHandler(suite.service)
	suite.router = gin.New()
	suite.router.POST("/seasons", suite.handler.Create)
	suite.router.GET("/seasons/:id", suite.handler.GetByID)
	suite.router.PUT("/seasons/:id", suite.handler.Update)
	suite.router.DELETE("/seasons/:id", suite.handler.Delete)
	suite.router.GET("/seasons", suite.handler.List)
}

func (suite *SeasonHandlerTestSuite) TestCreate() {
	req := CreateSeasonRequest{Title: "Test Season"}
	body, _ := json.Marshal(req)

	suite.repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Season")).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/seasons", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusCreated, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonHandlerTestSuite) TestCreate_InvalidJSON() {
	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/seasons", bytes.NewBuffer([]byte("invalid json")))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *SeasonHandlerTestSuite) TestCreate_MissingRequiredField() {
	req := CreateSeasonRequest{} // Missing required Title
	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/seasons", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *SeasonHandlerTestSuite) TestGetByID() {
	season := &models.Season{ID: 1, Title: "Test Season"}
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(season, nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/seasons/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	var resp models.Season
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal(season.Title, resp.Title)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonHandlerTestSuite) TestGetByID_InvalidID() {
	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/seasons/invalid", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *SeasonHandlerTestSuite) TestGetByID_NotFound() {
	suite.repo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/seasons/999", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonHandlerTestSuite) TestUpdate() {
	existingSeason := &models.Season{ID: 1, Title: "Old Title"}
	req := CreateSeasonRequest{Title: "Updated Title"}
	body, _ := json.Marshal(req)

	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(existingSeason, nil)
	suite.repo.On("Update", mock.Anything, mock.AnythingOfType("*models.Season")).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("PUT", "/seasons/1", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonHandlerTestSuite) TestUpdate_NotFound() {
	req := CreateSeasonRequest{Title: "Updated Title"}
	body, _ := json.Marshal(req)

	suite.repo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("PUT", "/seasons/999", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonHandlerTestSuite) TestDelete() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("DELETE", "/seasons/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonHandlerTestSuite) TestDelete_Error() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("DELETE", "/seasons/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *SeasonHandlerTestSuite) TestList() {
	seasons := []models.Season{
		{ID: 1, Title: "Season 1"},
		{ID: 2, Title: "Season 2"},
	}
	suite.repo.On("List", mock.Anything, 1, 10).Return(seasons, int64(2), nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/seasons?page=1&limit=10", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	var resp dto.ListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal(int64(2), resp.Total)
	suite.repo.AssertExpectations(suite.T())
}

func TestSeasonHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SeasonHandlerTestSuite))
}
