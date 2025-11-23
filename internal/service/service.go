package service

import (
	"context"

	"github.com/amit/gin-crud-demo/internal/models"
	"github.com/amit/gin-crud-demo/internal/repository"
)

type MovieService struct {
	repo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) Create(ctx context.Context, movie *models.Movie) error {
	return s.repo.Create(ctx, movie)
}

func (s *MovieService) GetByID(ctx context.Context, id int64) (*models.Movie, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MovieService) Update(ctx context.Context, movie *models.Movie) error {
	return s.repo.Update(ctx, movie)
}

func (s *MovieService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *MovieService) List(ctx context.Context, page, limit int) ([]models.Movie, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.List(ctx, page, limit)
}

type TVShowService struct {
	repo repository.TVShowRepository
}

func NewTVShowService(repo repository.TVShowRepository) *TVShowService {
	return &TVShowService{repo: repo}
}

func (s *TVShowService) Create(ctx context.Context, tvShow *models.TVShow) error {
	return s.repo.Create(ctx, tvShow)
}

func (s *TVShowService) GetByID(ctx context.Context, id int64) (*models.TVShow, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TVShowService) Update(ctx context.Context, tvShow *models.TVShow) error {
	return s.repo.Update(ctx, tvShow)
}

func (s *TVShowService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *TVShowService) List(ctx context.Context, page, limit int) ([]models.TVShow, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.List(ctx, page, limit)
}

type SeasonService struct {
	repo repository.SeasonRepository
}

func NewSeasonService(repo repository.SeasonRepository) *SeasonService {
	return &SeasonService{repo: repo}
}

func (s *SeasonService) Create(ctx context.Context, season *models.Season) error {
	return s.repo.Create(ctx, season)
}

func (s *SeasonService) GetByID(ctx context.Context, id int64) (*models.Season, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SeasonService) Update(ctx context.Context, season *models.Season) error {
	return s.repo.Update(ctx, season)
}

func (s *SeasonService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *SeasonService) List(ctx context.Context, page, limit int) ([]models.Season, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.List(ctx, page, limit)
}

type ViewSummaryService struct {
	repo repository.ViewSummaryRepository
}

func NewViewSummaryService(repo repository.ViewSummaryRepository) *ViewSummaryService {
	return &ViewSummaryService{repo: repo}
}

func (s *ViewSummaryService) Create(ctx context.Context, summary *models.ViewSummary) error {
	return s.repo.Create(ctx, summary)
}

func (s *ViewSummaryService) GetByID(ctx context.Context, id int64) (*models.ViewSummary, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ViewSummaryService) Update(ctx context.Context, summary *models.ViewSummary) error {
	return s.repo.Update(ctx, summary)
}

func (s *ViewSummaryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *ViewSummaryService) List(ctx context.Context, page, limit int) ([]models.ViewSummary, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.List(ctx, page, limit)
}
