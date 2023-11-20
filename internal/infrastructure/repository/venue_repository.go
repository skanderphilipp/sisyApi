package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/venue"
	"gorm.io/gorm"
)

type VenueRepository struct {
	db *gorm.DB
}

func NewVenueRepository(db *gorm.DB) *VenueRepository {
	return &VenueRepository{db: db}
}

func (repo *VenueRepository) FindAllByCursor(ctx context.Context, cursor string, limit int) ([]venue.Venue, string, error) {
	var venues []venue.Venue
	var nextCursor string

	// Assuming ID is used as the cursor
	query := repo.db.WithContext(ctx).Order("id ASC")

	if cursor != "" {
		query = query.Where("id > ?", cursor)
	}

	err := query.Limit(limit).Preload("Stages").Find(&venues).Error
	if err != nil {
		return nil, "", err
	}

	// Set next cursor
	if len(venues) > 0 {
		nextCursor = venues[len(venues)-1].ID.String() // Assuming ID is a UUID
	}

	return venues, nextCursor, nil
}

func (r *VenueRepository) Save(ctx context.Context, venue *venue.Venue) (*venue.Venue, error) {
	// Save the venue to the database
	result := r.db.WithContext(ctx).Save(venue)

	if result.Error != nil {
		return nil, fmt.Errorf("error saving venue: %v", result.Error)
	}

	return venue, nil
}

func (r *VenueRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	// Delete the venue with the given ID from the database
	result := r.db.WithContext(ctx).Delete(&venue.Venue{}, id)

	if result.Error != nil {
		return false, fmt.Errorf("error deleting venue: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("venue not found")
	}

	return true, nil
}

func (r *VenueRepository) Update(ctx context.Context, venue *venue.Venue) (*venue.Venue, error) {
	if err := r.db.WithContext(ctx).Save(venue).Error; err != nil {
		return nil, err
	}
	return venue, nil
}
