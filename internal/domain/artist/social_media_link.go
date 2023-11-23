package artist

import (
	"errors"
	"github.com/blnto/blnto_service/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// SocialMediaLink represents the social media model with soft delete
type SocialMediaLink struct {
	ID        uuid.UUID           `gorm:"type:uuid;primaryKey;" json:"id"`
	ArtistID  uuid.UUID           `gorm:"type:uuid;not null;" json:"-"`
	Platform  SocialMediaPlatform `gorm:"type:varchar(50);not null;" json:"platform"`
	Link      string              `gorm:"type:text;not null;" json:"link"`
	CreatedAt time.Time           `json:"-"`
	UpdatedAt time.Time           `json:"-"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"-"`
}

func (s *SocialMediaLink) ValidateLinkFormat() error {
	// Implement link format validation logic
	// For example, using regex to check if it's a valid URL
	if !utils.IsValidURL(s.Link) {
		return errors.New("invalid link format")
	}
	return nil
}

// BeforeCreate will set a UUID rather than numeric ID.
func (s *SocialMediaLink) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}
