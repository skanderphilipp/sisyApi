package artist

import (
	"errors"
	"time"

	"github.com/blnto/blnto_service/internal/infrastructure/api/artistApi"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Artist represents an artist record in the database
type Artist struct {
	ID               uuid.UUID         `gorm:"type:uuid;primaryKey;" json:"id"`
	Name             string            `gorm:"type:varchar(100);not null;" json:"name"`
	Location         string            `gorm:"type:varchar(100);" json:"location"`
	SCPromotedSet    string            `gorm:"type:text;column:sc_promoted_set" json:"soundcloudPromotedSet"`
	SocialMediaLinks []SocialMediaLink `gorm:"foreignKey:ArtistID" json:"SocialMediaLinks"`
	SCID             *int              `gorm:"column:sc_id;unique"`
	CreatedAt        time.Time         `json:"-"`
	UpdatedAt        time.Time         `json:"-"`
	DeletedAt        gorm.DeletedAt    `gorm:"index" json:"-"`
	SCCity           string            `gorm:"column:sc_city"`
	SCAvatarURL      string            `gorm:"column:sc_avatar_url"`
	SCFirstName      string            `gorm:"column:sc_first_name"`
	SCLastName       string            `gorm:"column:sc_last_name"`
	SCFullName       string            `gorm:"column:sc_full_name"`
	SCUsername       *string           `gorm:"column:sc_username;unique"`
	SCDescription    string            `gorm:"type:text;gorm:column:sc_description"`
	SCCountry        string            `gorm:"column:sc_country"`
	SCPermalink      *string           `gorm:"column:sc_permalink;unique"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (a *Artist) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
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

func (a *Artist) SetPermalink(permalink *string) error {
	if permalink != nil && *permalink == "" {
		return errors.New("permalink cannot be an empty string")
	}
	a.SCPermalink = permalink
	return nil
}
