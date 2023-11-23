package test

import (
	"github.com/blnto/blnto_service/internal/domain/artist"
)

// DirectCopy for testing
func DirectCopy(original *artist.Artist) *artist.Artist {
	copy := *original
	return &copy
}

// NewArtistCopy for testing
func NewArtistCopy(original *artist.Artist, modify func(*artist.Artist)) *artist.Artist {
	copy := *original
	if modify != nil {
		modify(&copy)
	}
	return &copy
}
