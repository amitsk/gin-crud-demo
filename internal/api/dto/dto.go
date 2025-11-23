package dto

import "time"

type CreateMovieRequest struct {
	Title             string  `json:"title" binding:"required"`
	OriginalTitle     *string `json:"original_title"`
	Runtime           *int64  `json:"runtime"`
	ReleaseDate       *string `json:"release_date"` // YYYY-MM-DD
	AvailableGlobally *bool   `json:"available_globally"`
	Locale            *string `json:"locale"`
}

type UpdateMovieRequest struct {
	Title             string  `json:"title"`
	OriginalTitle     *string `json:"original_title"`
	Runtime           *int64  `json:"runtime"`
	ReleaseDate       *string `json:"release_date"`
	AvailableGlobally *bool   `json:"available_globally"`
	Locale            *string `json:"locale"`
}

type MovieResponse struct {
	ID                int64     `json:"id"`
	Title             string    `json:"title"`
	OriginalTitle     *string   `json:"original_title"`
	Runtime           *int64    `json:"runtime"`
	ReleaseDate       *string   `json:"release_date"`
	AvailableGlobally *bool     `json:"available_globally"`
	Locale            *string   `json:"locale"`
	CreatedDate       time.Time `json:"created_date"`
	ModifiedDate      time.Time `json:"modified_date"`
}

type PaginationRequest struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}

type ListResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}
