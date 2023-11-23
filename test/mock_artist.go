package test

import (
	"github.com/blnto/blnto_service/internal/domain/artist"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func CreateMockArtist() *artist.Artist {
	// Sample data for SocialMediaLinks and SCID
	socialMediaLinks := []artist.SocialMediaLink{
		{
			// Populate with sample data
		},
		// Add more SocialMediaLink items as needed
	}

	sampleSCID := 123 // Sample SoundCloud ID
	username := gofakeit.Username()
	scPermalink := "http://example.com/permalink"
	return &artist.Artist{
		ID:               uuid.New(),
		Name:             "Sample Artist",
		Location:         "Sample Location",
		SCPromotedSet:    "Sample Promoted Set",
		SocialMediaLinks: socialMediaLinks,
		SCID:             &sampleSCID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        gorm.DeletedAt{},
		SCCity:           "Sample City",
		SCAvatarURL:      "http://example.com/avatar.jpg",
		SCFirstName:      "Sample",
		SCLastName:       "Artist",
		SCFullName:       "Sample Artist",
		SCUsername:       &username,
		SCDescription:    "Sample Description",
		SCCountry:        "Sample Country",
		SCPermalink:      &scPermalink,
	}
}
