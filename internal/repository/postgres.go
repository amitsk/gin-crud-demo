package repository

import (
	"context"

	"github.com/amit/gin-crud-demo/internal/models"
	"gorm.io/gorm"
)

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) Create(ctx context.Context, movie *models.Movie) error {
	return r.db.WithContext(ctx).Create(movie).Error
}

func (r *movieRepository) GetByID(ctx context.Context, id int64) (*models.Movie, error) {
	var movie models.Movie
	if err := r.db.WithContext(ctx).First(&movie, id).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *movieRepository) Update(ctx context.Context, movie *models.Movie) error {
	return r.db.WithContext(ctx).Save(movie).Error
}

func (r *movieRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Movie{}, id).Error
}

func (r *movieRepository) List(ctx context.Context, page, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Movie{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

type tvShowRepository struct {
	db *gorm.DB
}

func NewTVShowRepository(db *gorm.DB) TVShowRepository {
	return &tvShowRepository{db: db}
}

func (r *tvShowRepository) Create(ctx context.Context, tvShow *models.TVShow) error {
	return r.db.WithContext(ctx).Create(tvShow).Error
}

func (r *tvShowRepository) GetByID(ctx context.Context, id int64) (*models.TVShow, error) {
	var tvShow models.TVShow
	if err := r.db.WithContext(ctx).First(&tvShow, id).Error; err != nil {
		return nil, err
	}
	return &tvShow, nil
}

func (r *tvShowRepository) Update(ctx context.Context, tvShow *models.TVShow) error {
	return r.db.WithContext(ctx).Save(tvShow).Error
}

func (r *tvShowRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.TVShow{}, id).Error
}

func (r *tvShowRepository) List(ctx context.Context, page, limit int) ([]models.TVShow, int64, error) {
	var tvShows []models.TVShow
	var total int64
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.TVShow{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&tvShows).Error; err != nil {
		return nil, 0, err
	}

	return tvShows, total, nil
}

type seasonRepository struct {
	db *gorm.DB
}

func NewSeasonRepository(db *gorm.DB) SeasonRepository {
	return &seasonRepository{db: db}
}

func (r *seasonRepository) Create(ctx context.Context, season *models.Season) error {
	return r.db.WithContext(ctx).Create(season).Error
}

func (r *seasonRepository) GetByID(ctx context.Context, id int64) (*models.Season, error) {
	var season models.Season
	if err := r.db.WithContext(ctx).Preload("TVShow").First(&season, id).Error; err != nil {
		return nil, err
	}
	return &season, nil
}

func (r *seasonRepository) Update(ctx context.Context, season *models.Season) error {
	return r.db.WithContext(ctx).Save(season).Error
}

func (r *seasonRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Season{}, id).Error
}

func (r *seasonRepository) List(ctx context.Context, page, limit int) ([]models.Season, int64, error) {
	var seasons []models.Season
	var total int64
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Season{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("TVShow").Offset(offset).Limit(limit).Find(&seasons).Error; err != nil {
		return nil, 0, err
	}

	return seasons, total, nil
}

type viewSummaryRepository struct {
	db *gorm.DB
}

func NewViewSummaryRepository(db *gorm.DB) ViewSummaryRepository {
	return &viewSummaryRepository{db: db}
}

func (r *viewSummaryRepository) Create(ctx context.Context, summary *models.ViewSummary) error {
	return r.db.WithContext(ctx).Create(summary).Error
}

func (r *viewSummaryRepository) GetByID(ctx context.Context, id int64) (*models.ViewSummary, error) {
	var summary models.ViewSummary
	if err := r.db.WithContext(ctx).Preload("Movie").Preload("Season").First(&summary, id).Error; err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *viewSummaryRepository) Update(ctx context.Context, summary *models.ViewSummary) error {
	return r.db.WithContext(ctx).Save(summary).Error
}

func (r *viewSummaryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.ViewSummary{}, id).Error
}

func (r *viewSummaryRepository) List(ctx context.Context, page, limit int) ([]models.ViewSummary, int64, error) {
	var summaries []models.ViewSummary
	var total int64
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.ViewSummary{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Movie").Preload("Season").Offset(offset).Limit(limit).Find(&summaries).Error; err != nil {
		return nil, 0, err
	}

	return summaries, total, nil
}
