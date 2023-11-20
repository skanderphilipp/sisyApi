package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
	"gorm.io/gorm"
)

type ArtistRepository struct {
	db *gorm.DB
}

func NewArtistRepository(db *gorm.DB) *ArtistRepository {
	return &ArtistRepository{db: db}
}

// FindArtistsByCursor fetches a page of artists starting after the given cursor with the specified limit
func (repo *ArtistRepository) FindAllByCursor(ctx context.Context, cursor string, limit int) ([]artist.Artist, string, error) {
	var artists []artist.Artist
	var nextCursor string

	// Assuming ID is used as the cursor
	query := repo.db.WithContext(ctx).Order("id ASC")

	if cursor != "" {
		query = query.Where("id > ?", cursor)
	}

	err := query.Limit(limit).Preload("SocialMediaLinks").Find(&artists).Error
	if err != nil {
		return nil, "", err
	}

	// Set next cursor
	if len(artists) > 0 {
		nextCursor = artists[len(artists)-1].ID.String() // Assuming ID is a UUID
	}

	return artists, nextCursor, nil
}

func (repo *ArtistRepository) FindByID(ctx context.Context, id uuid.UUID) (*artist.Artist, error) {
	var artist artist.Artist
	if err := repo.db.WithContext(ctx).Preload("SocialMediaLinks").Where("id = ?", id).First(&artist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("artist not found")
		}
		return nil, err
	}
	return &artist, nil
}

func (repo *ArtistRepository) FindAllWithPermalink(ctx context.Context) ([]artist.Artist, error) {
	var artists []artist.Artist
	result := repo.db.WithContext(ctx).Where("sc_permalink IS NOT NULL AND sc_permalink != ''").Find(&artists)
	return artists, result.Error
}

func (repo *ArtistRepository) FindByName(ctx context.Context, name string) (*artist.Artist, error) {
	var artist artist.Artist
	if err := repo.db.WithContext(ctx).Preload("SocialMediaLinks").Where("name = ?", name).First(&artist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("artist not found")
		}
		return nil, err
	}
	return &artist, nil
}

func (r *ArtistRepository) Search(ctx context.Context, criteria *models.ArtistSearchInput) ([]artist.Artist, string, error) {
	// Create a GORM query builder
	query := r.db.Model(&artist.Artist{})
	searchTerm := strings.ToLower(*criteria.SearchTerm)
	var nextCursor string
	// Apply filters based on criteria
	if criteria != nil {
		if *criteria.SearchTerm != "" {
			// Build the query to search across all relevant columns.
			query = query.Where("lower(name) ILIKE ?", "%"+searchTerm+"%").
				Or("lower(location) ILIKE ?", "%"+searchTerm+"%").
				Or("lower(soundcloud_set_link) ILIKE ?", "%"+searchTerm+"%")
		}
	}

	// Execute the query and retrieve the matching records
	var artists []artist.Artist
	if err := query.Limit(*criteria.First).Find(&artists).Error; err != nil {
		return nil, "", err
	}
	// Set next cursor
	if len(artists) > 0 {
		nextCursor = artists[len(artists)-1].ID.String() // Assuming ID is a UUID
	}
	return artists, nextCursor, nil
}

func (r *ArtistRepository) Save(ctx context.Context, artist *artist.Artist) (*artist.Artist, error) {
	// Save the artist to the database
	result := r.db.WithContext(ctx).Save(artist)

	if result.Error != nil {
		return nil, fmt.Errorf("error saving artist: %v", result.Error)
	}

	return artist, nil
}

func (r *ArtistRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	// Delete the artist with the given ID from the database
	result := r.db.WithContext(ctx).Delete(&artist.Artist{}, id)

	if result.Error != nil {
		return false, fmt.Errorf("error deleting artist: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("artist not found")
	}

	return true, nil
}

func (r *ArtistRepository) Update(ctx context.Context, artist *artist.Artist) (*artist.Artist, error) {
	if err := r.db.WithContext(ctx).Save(artist).Error; err != nil {
		return nil, err
	}
	return artist, nil
}
