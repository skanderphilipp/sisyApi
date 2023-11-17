package repository

import (
	"context"

	"github.com/skanderphilipp/sisyApi/internal/domain/models"
	"gorm.io/gorm"
)

type VenueRepository struct {
	db *gorm.DB
}

func NewVenueRepository(db *gorm.DB) *VenueRepository {
	return &VenueRepository{db: db}
}

func (repo *VenueRepository) FindAllByCursor(ctx context.Context, cursor string, limit int) ([]models.Venue, string, error) {
	var venues []models.Venue
	var nextCursor string

	// Assuming ID is used as the cursor
	query := repo.db.WithContext(ctx).Order("id ASC")

	if cursor != "" {
		query = query.Where("id > ?", cursor)
	}

	err := query.Limit(limit).Preload("SocialMediaLinks").Find(&venues).Error
	if err != nil {
		return nil, "", err
	}

	// Set next cursor
	if len(venues) > 0 {
		nextCursor = venues[len(venues)-1].ID.String() // Assuming ID is a UUID
	}

	return venues, nextCursor, nil
}
