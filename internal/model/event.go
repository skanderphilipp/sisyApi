package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	UUID			uuid.UUID `gorm:"type:uuid;primaryKey"`
	VenueID		uuid.UUID `gorm:"type:uuid;foreignKey:VenueUUID"`
	startDate time.Time `gorm:"type:timestamp"`
	endDate time.Time `gorm:"type:timestamp"`
	Timetable []TimetableEntry `gorm:"foreignKey:EventID"`
}