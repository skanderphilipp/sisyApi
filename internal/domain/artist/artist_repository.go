package artist

import (
	"context"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
)

type Repository interface {
	// FindByID returns an artist by its ID
	FindByName(ctx context.Context, name string) (*models.Artist, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Artist, error)
	Search(ctx context.Context, criteria *models.ArtistSearchInput) ([]*models.Artist, string, error)
	FindAllByCursor(ctx context.Context, string, limit int) ([]*models.Artist, string, error)
	// SaveArtist saves an artist to the repository
	Save(ctx context.Context, artist *models.Artist) (*models.Artist, error)
	Update()
	// DeleteArtist deletes an artist from the repository by ID
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}
