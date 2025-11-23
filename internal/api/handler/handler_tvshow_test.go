package handler

import (
	"bytes"
	"encoding/json"
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

type TVShowHandlerTestSuite struct {
	suite.Suite
	repo    *mocks.MockTVShowRepository
	service *service.TVShowService
	handler *TVShowHandler
	router  *gin.Engine
}

func (suite *TVShowHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.repo = new(mocks.MockTVShowRepository)
	suite.service = service.NewTVShowService(suite.repo)
	suite.handler = NewTVShowHandler(suite.service)
	suite.router = gin.New()
	suite.router.POST("/tv-shows", suite.handler.Create)
	suite.router.GET("/tv-shows/:id", suite.handler.GetByID)
	suite.router.PUT("/tv-shows/:id", suite.handler.Update)
	suite.router.DELETE("/tv-shows/:id", suite.handler.Delete)
	suite.router.GET("/tv-shows", suite.handler.List)
}

func (suite *TVShowHandlerTestSuite) TestCreate() {
	req := dto.CreateMovieRequest{Title: "Test TV Show"}
	body, _ := json.Marshal(req)

	suite.repo.On("Create", mock.Anything, mock.AnythingOfType("*models.TVShow")).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/tv-shows", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusCreated, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowHandlerTestSuite) TestCreate_InvalidJSON() {
	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/tv-shows", bytes.NewBuffer([]byte("invalid json")))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *TVShowHandlerTestSuite) TestCreate_MissingRequiredField() {
	req := dto.CreateMovieRequest{} // Missing required Title
	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/tv-shows", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *TVShowHandlerTestSuite) TestGetByID() {
	tvShow := &models.TVShow{ID: 1, Title: "Test TV Show"}
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(tvShow, nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/tv-shows/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	var resp models.TVShow
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal(tvShow.Title, resp.Title)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowHandlerTestSuite) TestGetByID_InvalidID() {
	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/tv-shows/invalid", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *TVShowHandlerTestSuite) TestGetByID_NotFound() {
	suite.repo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/tv-shows/999", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowHandlerTestSuite) TestUpdate() {
	existingTVShow := &models.TVShow{ID: 1, Title: "Old Title"}
	req := dto.UpdateMovieRequest{Title: "Updated Title"}
	body, _ := json.Marshal(req)

	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(existingTVShow, nil)
	suite.repo.On("Update", mock.Anything, mock.AnythingOfType("*models.TVShow")).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("PUT", "/tv-shows/1", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowHandlerTestSuite) TestUpdate_NotFound() {
	req := dto.UpdateMovieRequest{Title: "Updated Title"}
	body, _ := json.Marshal(req)

	suite.repo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("PUT", "/tv-shows/999", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowHandlerTestSuite) TestDelete() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("DELETE", "/tv-shows/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowHandlerTestSuite) TestDelete_Error() {
	suite.repo.On("Delete", mock.Anything, int64(1)).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("DELETE", "/tv-shows/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TVShowHandlerTestSuite) TestList() {
	tvShows := []models.TVShow{
		{ID: 1, Title: "TV Show 1"},
		{ID: 2, Title: "TV Show 2"},
	}
	suite.repo.On("List", mock.Anything, 1, 10).Return(tvShows, int64(2), nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/tv-shows?page=1&limit=10", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	var resp dto.ListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal(int64(2), resp.Total)
	suite.repo.AssertExpectations(suite.T())
}

func TestTVShowHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TVShowHandlerTestSuite))
}
