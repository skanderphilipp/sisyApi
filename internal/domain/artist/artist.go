package artist

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Artist represents an artist record in the database
type Artist struct {
	UUID              uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	Name              string         `gorm:"type:varchar(100);not null;" json:"name"`
	Location          string         `gorm:"type:varchar(100);" json:"location"`
	SoundcloudSetLink string         `gorm:"type:text;" json:"soundcloudSetLink"`
	SocialMediaLinks  []SocialMedia  `gorm:"foreignKey:ArtistID" json:"SocialMediaLinks"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type SocialMediaPlatform string

const (
	Twitter    SocialMediaPlatform = "Twitter"
	Facebook   SocialMediaPlatform = "Facebook"
	Instagram  SocialMediaPlatform = "Instagram"
	YouTube    SocialMediaPlatform = "YouTube"
	Soundcloud SocialMediaPlatform = "Soundcloud"
)

// SocialMedia represents the social media model with soft delete
type SocialMedia struct {
	UUID                uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	ArtistID            uuid.UUID      `gorm:"type:uuid;not null;" json:"-"`
	SocialMediaPlatform string         `gorm:"type:varchar(50);not null;" json:"platform"`
	Link                string         `gorm:"type:text;not null;" json:"link"`
	CreatedAt           time.Time      `json:"-"`
	UpdatedAt           time.Time      `json:"-"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (artist *Artist) BeforeCreate(tx *gorm.DB) (err error) {
	artist.UUID = uuid.New()
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (socialMedia *SocialMedia) BeforeCreate(tx *gorm.DB) (err error) {
	socialMedia.UUID = uuid.New()
	return
}
