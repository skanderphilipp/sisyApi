package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
)

type TimetableEntry struct {
	UUID       uuid.UUID     `gorm:"type:uuid;primaryKey"`
	EventID    uuid.UUID     `gorm:"type:uuid"`
	StageID    uuid.UUID     `gorm:"type:uuid"`
	Stage      Stage         `gorm:"foreignKey:StageID"`
	ArtistID   uuid.UUID     `gorm:"type:uuid"`
	Artist     artist.Artist `gorm:"foreignKey:ArtistID"`
	WeekNumber int
	Year       int
	Day        string `gorm:"type:varchar"`
	StartTime  time.Time
	EndTime    time.Time
}
