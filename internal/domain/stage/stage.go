package stage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Stage represents a stage in a venue.
type Stage struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StageName string         `gorm:"type:varchar(100);not null" json:"stageName"`
	VenueID   uuid.UUID      `gorm:"type:uuid;foreignKey:VenueID" json:"venueID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Stage BeforeCreate hook
func (s *Stage) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
