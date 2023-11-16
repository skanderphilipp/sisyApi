package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
	"gorm.io/gorm"
)

type ArtistRepository struct {
	db *gorm.DB
}

func NewArtistRepository(db *gorm.DB) *ArtistRepository {
	return &ArtistRepository{db: db}
}

func (r *ArtistRepository) GetByID(ctx context.Context, id uuid.UUID) (*artist.Artist, error) {
	var artist artist.Artist
	if err := r.db.WithContext(ctx).First(&artist, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}

// Implement other methods...
