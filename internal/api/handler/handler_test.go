package handler

import (
	"bytes"
	"encoding/json"
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

type MovieHandlerTestSuite struct {
	suite.Suite
	repo    *mocks.MockMovieRepository
	service *service.MovieService
	handler *MovieHandler
	router  *gin.Engine
}

func (suite *MovieHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.repo = new(mocks.MockMovieRepository)
	suite.service = service.NewMovieService(suite.repo)
	suite.handler = NewMovieHandler(suite.service)
	suite.router = gin.New()
	suite.router.POST("/movies", suite.handler.Create)
	suite.router.GET("/movies/:id", suite.handler.GetByID)
}

func (suite *MovieHandlerTestSuite) TestCreate() {
	req := dto.CreateMovieRequest{Title: "Test Movie"}
	body, _ := json.Marshal(req)

	suite.repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Movie")).Return(nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("POST", "/movies", bytes.NewBuffer(body))
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusCreated, w.Code)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *MovieHandlerTestSuite) TestGetByID() {
	movie := &models.Movie{ID: 1, Title: "Test Movie"}
	suite.repo.On("GetByID", mock.Anything, int64(1)).Return(movie, nil)

	w := httptest.NewRecorder()
	reqHTTP, _ := http.NewRequest("GET", "/movies/1", nil)
	suite.router.ServeHTTP(w, reqHTTP)

	suite.Equal(http.StatusOK, w.Code)
	var resp models.Movie
	json.Unmarshal(w.Body.Bytes(), &resp)
	suite.Equal(movie.Title, resp.Title)
	suite.repo.AssertExpectations(suite.T())
}

func TestMovieHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(MovieHandlerTestSuite))
}
