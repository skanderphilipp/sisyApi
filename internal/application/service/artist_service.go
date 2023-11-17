package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
	"github.com/skanderphilipp/sisyApi/internal/infrastructure/repository"
)

type ArtistService struct {
	repo *repository.ArtistRepository
}

func NewArtistService(repo *repository.ArtistRepository) *ArtistService {
	return &ArtistService{repo: repo}
}

func (s *ArtistService) GetArtist(ctx context.Context, id uuid.UUID) (*models.Artist, error) {
	return s.repo.FindByID(ctx, id)
}
