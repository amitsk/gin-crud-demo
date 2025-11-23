package models

import (
	"time"
)

type Season struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedDate   time.Time `gorm:"not null" json:"created_date"`
	ModifiedDate  time.Time `gorm:"not null" json:"modified_date"`
	OriginalTitle *string   `gorm:"type:varchar(255)" json:"original_title"`
	ReleaseDate   *time.Time `gorm:"type:date" json:"release_date"`
	Runtime       *int64    `json:"runtime"`
	SeasonNumber  *int      `json:"season_number"`
	Title         string    `gorm:"type:varchar(255);not null" json:"title"`
	TVShowID      *int64    `json:"tv_show_id"`
	TVShow        *TVShow   `gorm:"foreignKey:TVShowID" json:"tv_show,omitempty"`
}

func (Season) TableName() string {
	return "season"
}
