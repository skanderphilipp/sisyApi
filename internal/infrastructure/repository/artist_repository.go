package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/blnto/blnto_service/internal/domain/artist"
	"github.com/blnto/blnto_service/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArtistRepository struct {
	db *gorm.DB
}

func NewArtistRepository(db *gorm.DB) *ArtistRepository {
	return &ArtistRepository{db: db}
}

// FindAllByCursor FindArtistsByCursor fetches a page of artists starting after the given cursor with the specified limit
func (r *ArtistRepository) FindAllByCursor(ctx context.Context, cursor string, limit int) ([]artist.Artist, string, error) {
	var artists []artist.Artist
	var nextCursor string

	// Assuming ID is used as the cursor
	query := r.db.WithContext(ctx).Order("id ASC")

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

func (r *ArtistRepository) FindByID(ctx context.Context, id uuid.UUID) (*artist.Artist, error) {
	var (
		artistModel artist.Artist
	)
	if err := r.db.WithContext(ctx).Preload("SocialMediaLinks").Where("id = ?", id).First(&artistModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("artistModel not found")
		}
		return nil, err
	}
	return &artistModel, nil
}

func (r *ArtistRepository) FindAllWithPermalink(ctx context.Context) ([]artist.Artist, error) {
	var artists []artist.Artist
	result := r.db.WithContext(ctx).Where("sc_permalink IS NOT NULL AND sc_permalink != ''").Find(&artists)
	return artists, result.Error
}

func (r *ArtistRepository) FindByName(ctx context.Context, name string) (*artist.Artist, error) {
	var (
		artistModel artist.Artist
	)
	if err := r.db.WithContext(ctx).Preload("SocialMediaLinks").Where("name = ?", name).First(&artistModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("artistModel not found")
		}
		return nil, err
	}
	return &artistModel, nil
}

func (r *ArtistRepository) Search(ctx context.Context, criteria *models.ArtistSearchInput) ([]artist.Artist, string, error) {
	// Create a GORM query builder
	query := r.db.WithContext(ctx).Model(&artist.Artist{})
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

func (r *ArtistRepository) PermalinkExists(ctx context.Context, permalink string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&artist.Artist{}).Where("sc_permalink = ?", permalink).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("error checking if permalink exists: %v", result.Error)
	}
	return count > 0, nil
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

func (r *ArtistRepository) CreateSocialMediaLink(ctx context.Context, link artist.SocialMediaLink) error {
	db := r.db.WithContext(ctx)
	return db.Create(&link).Error
}

func (r *ArtistRepository) UpdateSocialMediaLink(ctx context.Context, link artist.SocialMediaLink) error {
	db := r.db.WithContext(ctx)
	return db.Save(&link).Error
}

func (r *ArtistRepository) GetSocialMediaLinksByArtistID(ctx context.Context, id uuid.UUID) ([]*artist.SocialMediaLink, error) {
	var links []*artist.SocialMediaLink
	db := r.db.WithContext(ctx)
	result := db.Where("artist_id = ?", id).Find(&links)
	return links, result.Error
}

func (r *ArtistRepository) DeleteSocialMediaLink(ctx context.Context, id uuid.UUID) error {
	db := r.db.WithContext(ctx)
	return db.Delete(&artist.SocialMediaLink{}, "id = ?", id).Error
}

func (r *ArtistRepository) PermalinkExistsExcludingArtist(ctx context.Context, permalink string, id uuid.UUID) (bool, error) {
	var (
		count int64
	)
	db := r.db.WithContext(ctx)
	result := db.Model(&artist.Artist{}).
		Where("sc_permalink = ?", permalink).
		Where("id <> ?", id).
		Count(&count)

	return count > 0, result.Error
}

func (r *ArtistRepository) FindFeatured(ctx context.Context) ([]artist.Artist, error) {
	var artists []artist.Artist
	db := r.db.WithContext(ctx).
		Model(&artist.Artist{}).
		Where("sc_id IS NOT NULL").
		Preload("SocialMediaLinks").
		Limit(10).Find(&artists)
	if db.Error != nil {
		_ = fmt.Errorf("error checking for featured Artists: %v", db.Error)
		return nil, db.Error
	}
	return artists, nil
}
