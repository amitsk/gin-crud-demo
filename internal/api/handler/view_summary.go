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

type ViewSummaryHandler struct {
	service *service.ViewSummaryService
}

func NewViewSummaryHandler(service *service.ViewSummaryService) *ViewSummaryHandler {
	return &ViewSummaryHandler{service: service}
}

type CreateViewSummaryRequest struct {
	CumulativeWeeksInTop10 *int   `json:"cumulative_weeks_in_top10"`
	Duration               string `json:"duration" binding:"required"`
	EndDate                string `json:"end_date" binding:"required"`
	HoursViewed            int    `json:"hours_viewed" binding:"required"`
	StartDate              string `json:"start_date" binding:"required"`
	ViewRank               *int   `json:"view_rank"`
	Views                  *int   `json:"views"`
	MovieID                *int64 `json:"movie_id"`
	SeasonID               *int64 `json:"season_id"`
}

func (h *ViewSummaryHandler) Create(c *gin.Context) {
	var req CreateViewSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	summary := &models.ViewSummary{
		CumulativeWeeksInTop10: req.CumulativeWeeksInTop10,
		Duration:               req.Duration,
		EndDate:                endDate,
		HoursViewed:            req.HoursViewed,
		StartDate:              startDate,
		ViewRank:               req.ViewRank,
		Views:                  req.Views,
		MovieID:                req.MovieID,
		SeasonID:               req.SeasonID,
		CreatedDate:            time.Now(),
		ModifiedDate:           time.Now(),
	}

	if err := h.service.Create(c.Request.Context(), summary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, summary)
}

func (h *ViewSummaryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	summary, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "view summary not found"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *ViewSummaryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req CreateViewSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "view summary not found"})
		return
	}

	if req.CumulativeWeeksInTop10 != nil {
		summary.CumulativeWeeksInTop10 = req.CumulativeWeeksInTop10
	}
	if req.Duration != "" {
		summary.Duration = req.Duration
	}
	if req.HoursViewed != 0 {
		summary.HoursViewed = req.HoursViewed
	}
	if req.ViewRank != nil {
		summary.ViewRank = req.ViewRank
	}
	if req.Views != nil {
		summary.Views = req.Views
	}
	if req.MovieID != nil {
		summary.MovieID = req.MovieID
	}
	if req.SeasonID != nil {
		summary.SeasonID = req.SeasonID
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			summary.StartDate = startDate
		}
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			summary.EndDate = endDate
		}
	}

	summary.ModifiedDate = time.Now()

	if err := h.service.Update(c.Request.Context(), summary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *ViewSummaryHandler) Delete(c *gin.Context) {
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

func (h *ViewSummaryHandler) List(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summaries, total, err := h.service.List(c.Request.Context(), req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{
		Data:  summaries,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	})
}
