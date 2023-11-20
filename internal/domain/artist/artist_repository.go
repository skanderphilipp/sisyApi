package artist

import (
	"context"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
)

type Repository interface {
	// FindByID returns an artist by its ID
	FindByName(ctx context.Context, name string) (*Artist, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Artist, error)
	FindByPermalink(ctx context.Context, permalink string) (*Artist, error)
	Search(ctx context.Context, criteria *models.ArtistSearchInput) ([]*Artist, string, error)
	FindAllByCursor(ctx context.Context, string, limit int) ([]*Artist, string, error)
	// SaveArtist saves an artist to the repository
	Save(ctx context.Context, artist *Artist) (*Artist, error)
	// DeleteArtist deletes an artist from the repository by ID
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	Update(ctx context.Context, artist *Artist) (*Artist, error)
}
