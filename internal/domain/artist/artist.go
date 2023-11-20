package artist

import (
	"time"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/infrastructure/api/artistApi"
	"gorm.io/gorm"
)

// Artist represents an artist record in the database
type Artist struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	Name             string         `gorm:"type:varchar(100);not null;" json:"name"`
	Location         string         `gorm:"type:varchar(100);" json:"location"`
	SCPromotedSet    string         `gorm:"type:text;column:sc_promoted_set" json:"soundcloudPromotedSet"`
	SocialMediaLinks []SocialMedia  `gorm:"foreignKey:ArtistID" json:"SocialMediaLinks"`
	SCID             *int           `gorm:"column:sc_id;unique"`
	CreatedAt        time.Time      `json:"-"`
	UpdatedAt        time.Time      `json:"-"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	SCCity           string         `gorm:"column:sc_city"`
	SCAvatarURL      string         `gorm:"column:sc_avatar_url"`
	SCFirstName      string         `gorm:"column:sc_first_name"`
	SCLastName       string         `gorm:"column:sc_last_name"`
	SCFullName       string         `gorm:"column:sc_full_name"`
	SCUsername       *string        `gorm:"column:sc_username;unique"`
	SCDescription    string         `gorm:"type:text;gorm:column:sc_description"`
	SCCountry        string         `gorm:"column:sc_country"`
	SCPermalink      *string        `gorm:"column:sc_permalink;unique"`
}

// SocialMedia represents the social media model with soft delete
type SocialMedia struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	ArtistID  uuid.UUID      `gorm:"type:uuid;not null;" json:"-"`
	Platform  string         `gorm:"type:varchar(50);not null;" json:"platform"`
	Link      string         `gorm:"type:text;not null;" json:"link"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (artist *Artist) BeforeCreate(tx *gorm.DB) (err error) {
	artist.ID = uuid.New()
	return
}

// BeforeCreate will set a UUID rather than numeric ID.
func (socialMedia *SocialMedia) BeforeCreate(tx *gorm.DB) (err error) {
	socialMedia.ID = uuid.New()
	return
}

func (a *Artist) UpdateWithSoundCloudInfo(scArtist artistApi.SCArtist) {
	a.SCUsername = &scArtist.Username
	a.SCFullName = scArtist.FullName
	a.SCAvatarURL = scArtist.AvatarURL
	a.SCCity = scArtist.City
	a.SCCountry = scArtist.Country
	a.SCDescription = scArtist.Description
	a.SCID = &scArtist.ID
}
