package repository

import (
	"context"

	"github.com/amit/gin-crud-demo/internal/models"
)

type MovieRepository interface {
	Create(ctx context.Context, movie *models.Movie) error
	GetByID(ctx context.Context, id int64) (*models.Movie, error)
	Update(ctx context.Context, movie *models.Movie) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int) ([]models.Movie, int64, error)
}

type TVShowRepository interface {
	Create(ctx context.Context, tvShow *models.TVShow) error
	GetByID(ctx context.Context, id int64) (*models.TVShow, error)
	Update(ctx context.Context, tvShow *models.TVShow) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int) ([]models.TVShow, int64, error)
}

type SeasonRepository interface {
	Create(ctx context.Context, season *models.Season) error
	GetByID(ctx context.Context, id int64) (*models.Season, error)
	Update(ctx context.Context, season *models.Season) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int) ([]models.Season, int64, error)
}

type ViewSummaryRepository interface {
	Create(ctx context.Context, summary *models.ViewSummary) error
	GetByID(ctx context.Context, id int64) (*models.ViewSummary, error)
	Update(ctx context.Context, summary *models.ViewSummary) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int) ([]models.ViewSummary, int64, error)
}
