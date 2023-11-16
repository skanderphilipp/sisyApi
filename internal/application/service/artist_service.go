package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
)

type ArtistService struct {
	repo artist.Repository
}

func NewArtistService(repo artist.Repository) *ArtistService {
	return &ArtistService{repo: repo}
}

func (s *ArtistService) GetArtist(ctx context.Context, id uuid.UUID) (*artist.Artist, error) {
	return s.repo.GetByID(ctx, id)
}
