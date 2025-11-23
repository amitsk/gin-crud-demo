package api

import (
	"github.com/amit/gin-crud-demo/internal/api/handler"
	"github.com/amit/gin-crud-demo/internal/middleware"
	"github.com/amit/gin-crud-demo/internal/repository"
	"github.com/amit/gin-crud-demo/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()

	// Setup Middleware
	middleware.Setup(r)

	// Repositories
	movieRepo := repository.NewMovieRepository(db)
	tvShowRepo := repository.NewTVShowRepository(db)
	seasonRepo := repository.NewSeasonRepository(db)
	viewSummaryRepo := repository.NewViewSummaryRepository(db)

	// Services
	movieService := service.NewMovieService(movieRepo)
	tvShowService := service.NewTVShowService(tvShowRepo)
	seasonService := service.NewSeasonService(seasonRepo)
	viewSummaryService := service.NewViewSummaryService(viewSummaryRepo)

	// Handlers
	movieHandler := handler.NewMovieHandler(movieService)
	tvShowHandler := handler.NewTVShowHandler(tvShowService)
	seasonHandler := handler.NewSeasonHandler(seasonService)
	viewSummaryHandler := handler.NewViewSummaryHandler(viewSummaryService)

	// Routes
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/healthz", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Movies
		movies := v1.Group("/movies")
		{
			movies.POST("", movieHandler.Create)
			movies.GET("/:id", movieHandler.GetByID)
			movies.PUT("/:id", movieHandler.Update)
			movies.DELETE("/:id", movieHandler.Delete)
			movies.GET("", movieHandler.List)
		}

		// TV Shows
		tvShows := v1.Group("/tv-shows")
		{
			tvShows.POST("", tvShowHandler.Create)
			tvShows.GET("/:id", tvShowHandler.GetByID)
			tvShows.PUT("/:id", tvShowHandler.Update)
			tvShows.DELETE("/:id", tvShowHandler.Delete)
			tvShows.GET("", tvShowHandler.List)
		}

		// Seasons
		seasons := v1.Group("/seasons")
		{
			seasons.POST("", seasonHandler.Create)
			seasons.GET("/:id", seasonHandler.GetByID)
			seasons.PUT("/:id", seasonHandler.Update)
			seasons.DELETE("/:id", seasonHandler.Delete)
			seasons.GET("", seasonHandler.List)
		}

		// View Summaries
		viewSummaries := v1.Group("/view-summaries")
		{
			viewSummaries.POST("", viewSummaryHandler.Create)
			viewSummaries.GET("/:id", viewSummaryHandler.GetByID)
			viewSummaries.PUT("/:id", viewSummaryHandler.Update)
			viewSummaries.DELETE("/:id", viewSummaryHandler.Delete)
			viewSummaries.GET("", viewSummaryHandler.List)
		}
	}

	return r
}
