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

type TVShowHandler struct {
	service *service.TVShowService
}

func NewTVShowHandler(service *service.TVShowService) *TVShowHandler {
	return &TVShowHandler{service: service}
}

func (h *TVShowHandler) Create(c *gin.Context) {
	var req dto.CreateMovieRequest // Reusing similar request struct for simplicity, ideally should have specific DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tvShow := &models.TVShow{
		Title:             req.Title,
		OriginalTitle:     req.OriginalTitle,
		AvailableGlobally: req.AvailableGlobally,
		Locale:            req.Locale,
		CreatedDate:       time.Now(),
		ModifiedDate:      time.Now(),
	}

	if req.ReleaseDate != nil {
		date, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err == nil {
			tvShow.ReleaseDate = &date
		}
	}

	if err := h.service.Create(c.Request.Context(), tvShow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tvShow)
}

func (h *TVShowHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tvShow, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tv show not found"})
		return
	}

	c.JSON(http.StatusOK, tvShow)
}

func (h *TVShowHandler) Update(c *gin.Context) {
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

	tvShow, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tv show not found"})
		return
	}

	if req.Title != "" {
		tvShow.Title = req.Title
	}
	if req.OriginalTitle != nil {
		tvShow.OriginalTitle = req.OriginalTitle
	}
	if req.AvailableGlobally != nil {
		tvShow.AvailableGlobally = req.AvailableGlobally
	}
	if req.Locale != nil {
		tvShow.Locale = req.Locale
	}
	if req.ReleaseDate != nil {
		date, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err == nil {
			tvShow.ReleaseDate = &date
		}
	}
	tvShow.ModifiedDate = time.Now()

	if err := h.service.Update(c.Request.Context(), tvShow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tvShow)
}

func (h *TVShowHandler) Delete(c *gin.Context) {
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

func (h *TVShowHandler) List(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tvShows, total, err := h.service.List(c.Request.Context(), req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{
		Data:  tvShows,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	})
}
