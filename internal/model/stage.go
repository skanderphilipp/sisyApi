package model

import "github.com/google/uuid"

type Stage struct {
	UUID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	StageName string    `gorm:"type:varchar"`
	VenueID uuid.UUID `gorm:"type:uuid;foreignKey:VenueUUID"`
}
