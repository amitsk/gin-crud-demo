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

type SeasonHandler struct {
	service *service.SeasonService
}

func NewSeasonHandler(service *service.SeasonService) *SeasonHandler {
	return &SeasonHandler{service: service}
}

type CreateSeasonRequest struct {
	Title         string  `json:"title" binding:"required"`
	OriginalTitle *string `json:"original_title"`
	Runtime       *int64  `json:"runtime"`
	ReleaseDate   *string `json:"release_date"`
	SeasonNumber  *int    `json:"season_number"`
	TVShowID      *int64  `json:"tv_show_id"`
}

func (h *SeasonHandler) Create(c *gin.Context) {
	var req CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	season := &models.Season{
		Title:         req.Title,
		OriginalTitle: req.OriginalTitle,
		Runtime:       req.Runtime,
		SeasonNumber:  req.SeasonNumber,
		TVShowID:      req.TVShowID,
		CreatedDate:   time.Now(),
		ModifiedDate:  time.Now(),
	}

	if req.ReleaseDate != nil {
		date, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err == nil {
			season.ReleaseDate = &date
		}
	}

	if err := h.service.Create(c.Request.Context(), season); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, season)
}

func (h *SeasonHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	season, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "season not found"})
		return
	}

	c.JSON(http.StatusOK, season)
}

func (h *SeasonHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	season, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "season not found"})
		return
	}

	if req.Title != "" {
		season.Title = req.Title
	}
	if req.OriginalTitle != nil {
		season.OriginalTitle = req.OriginalTitle
	}
	if req.Runtime != nil {
		season.Runtime = req.Runtime
	}
	if req.SeasonNumber != nil {
		season.SeasonNumber = req.SeasonNumber
	}
	if req.TVShowID != nil {
		season.TVShowID = req.TVShowID
	}
	if req.ReleaseDate != nil {
		date, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err == nil {
			season.ReleaseDate = &date
		}
	}
	season.ModifiedDate = time.Now()

	if err := h.service.Update(c.Request.Context(), season); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, season)
}

func (h *SeasonHandler) Delete(c *gin.Context) {
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

func (h *SeasonHandler) List(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seasons, total, err := h.service.List(c.Request.Context(), req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{
		Data:  seasons,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	})
}
