package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/amit/gin-crud-demo/internal/api/dto"
	"github.com/amit/gin-crud-demo/internal/models"
	"github.com/amit/gin-crud-demo/internal/service"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	service *service.MovieService
}

func NewMovieHandler(service *service.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (h *MovieHandler) Create(c *gin.Context) {
	var req dto.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie := &models.Movie{
		Title:             req.Title,
		OriginalTitle:     req.OriginalTitle,
		Runtime:           req.Runtime,
		AvailableGlobally: req.AvailableGlobally,
		Locale:            req.Locale,
		CreatedDate:       time.Now(),
		ModifiedDate:      time.Now(),
	}

	if req.ReleaseDate != nil {
		date, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err == nil {
			movie.ReleaseDate = &date
		}
	}

	if err := h.service.Create(c.Request.Context(), movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	movie, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	if req.Title != "" {
		movie.Title = req.Title
	}
	if req.OriginalTitle != nil {
		movie.OriginalTitle = req.OriginalTitle
	}
	if req.Runtime != nil {
		movie.Runtime = req.Runtime
	}
	if req.AvailableGlobally != nil {
		movie.AvailableGlobally = req.AvailableGlobally
	}
	if req.Locale != nil {
		movie.Locale = req.Locale
	}
	if req.ReleaseDate != nil {
		date, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err == nil {
			movie.ReleaseDate = &date
		}
	}
	movie.ModifiedDate = time.Now()

	if err := h.service.Update(c.Request.Context(), movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *MovieHandler) List(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movies, total, err := h.service.List(c.Request.Context(), req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{
		Data:  movies,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	})
}
