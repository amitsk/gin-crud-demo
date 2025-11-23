package models

import (
	"time"
)

type ViewSummary struct {
	ID                     int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedDate            time.Time `gorm:"not null" json:"created_date"`
	ModifiedDate           time.Time `gorm:"not null" json:"modified_date"`
	CumulativeWeeksInTop10 *int      `json:"cumulative_weeks_in_top10"`
	Duration               string    `gorm:"type:varchar(20);not null;check:duration IN ('WEEKLY','SEMI_ANNUALLY')" json:"duration"`
	EndDate                time.Time `gorm:"type:date;not null" json:"end_date"`
	HoursViewed            int       `gorm:"not null" json:"hours_viewed"`
	StartDate              time.Time `gorm:"type:date;not null" json:"start_date"`
	ViewRank               *int      `json:"view_rank"`
	Views                  *int      `json:"views"`
	MovieID                *int64    `json:"movie_id"`
	Movie                  *Movie    `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	SeasonID               *int64    `json:"season_id"`
	Season                 *Season   `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
}

func (ViewSummary) TableName() string {
	return "view_summary"
}
