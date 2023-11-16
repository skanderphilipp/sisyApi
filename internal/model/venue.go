package model

import "github.com/google/uuid"


type Venue struct {
	UUID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string `gorm:"type:varchar"`
	Description string `gorm:"type:varchar"`
	Stages []Stage `gorm:"foreignKey:VenueUUID"`
}