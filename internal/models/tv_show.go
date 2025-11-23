package models

import (
	"time"
)

type TVShow struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedDate       time.Time `gorm:"not null" json:"created_date"`
	ModifiedDate      time.Time `gorm:"not null" json:"modified_date"`
	AvailableGlobally *bool     `json:"available_globally"`
	Locale            *string   `gorm:"type:varchar(10)" json:"locale"`
	OriginalTitle     *string   `gorm:"type:varchar(255)" json:"original_title"`
	ReleaseDate       *time.Time `gorm:"type:date" json:"release_date"`
	Title             string    `gorm:"type:varchar(255);not null;unique" json:"title"`
}

func (TVShow) TableName() string {
	return "tv_show"
}
