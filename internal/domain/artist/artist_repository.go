package artist

import (
	"context"

	"github.com/blnto/blnto_service/internal/domain/models"
	"github.com/google/uuid"
)

type Repository interface {
	// FindByName returns an artist by its ID
	FindByName(ctx context.Context, name string) (*Artist, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Artist, error)
	FindByPermalink(ctx context.Context, permalink string) (*Artist, error)
	PermalinkExists(ctx context.Context, permalink string) (bool, error)
	Search(ctx context.Context, criteria *models.ArtistSearchInput) ([]*Artist, string, error)
	FindAllByCursor(ctx context.Context, string, limit int) ([]*Artist, string, error)
	// Save SaveArtist saves an artist to the repository
	Save(ctx context.Context, artist *Artist) (*Artist, error)
	// Delete DeleteArtist deletes an artist from the repository by ID
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	Update(ctx context.Context, artist *Artist) (*Artist, error)
	GetSocialMediaLinksByArtistID(ctx context.Context, artistID uuid.UUID) ([]*SocialMediaLink, error)
}
