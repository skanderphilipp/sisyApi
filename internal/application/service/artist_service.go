package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
	"github.com/skanderphilipp/sisyApi/internal/infrastructure/repository"
)

type ArtistService struct {
	repo *repository.ArtistRepository
}

func NewArtistService(repo *repository.ArtistRepository) *ArtistService {
	return &ArtistService{repo: repo}
}

func (s *ArtistService) GetArtist(ctx context.Context, id uuid.UUID) (*artist.Artist, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ArtistService) UpdateArtist(ctx context.Context, gqlArtist *models.Artist) (*models.Artist, error) {
	// Map GraphQL artist to GORM artist
	artistModel := mapGqlArtistToGormArtist(gqlArtist)

	// Update artist in the database
	updatedArtist, err := s.repo.Update(ctx, artistModel)
	if err != nil {
		return nil, err
	}

	// Map updated GORM artist back to GraphQL artist
	return mapGormArtistToGqlArtist(updatedArtist), nil
}

func (s *ArtistService) FindByName(ctx context.Context, name string) (*models.Artist, error) {
	artist, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return mapGormArtistToGqlArtist(artist), nil
}

func (s *ArtistService) FindByID(ctx context.Context, id uuid.UUID) (*models.Artist, error) {
	artist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapGormArtistToGqlArtist(artist), nil
}

func (s *ArtistService) Search(ctx context.Context, criteria *models.ArtistSearchInput) ([]*models.Artist, string, error) {
	// Implement your search logic
	artists, nextCursor, err := s.repo.Search(ctx, criteria)

	if err != nil {
		return nil, "", err
	}

	var wg sync.WaitGroup
	result := make([]*models.Artist, len(artists))
	errs := make(chan error, 1) // Buffered channel to avoid blocking

	for i, gormArtist := range artists {
		wg.Add(1) // Increment the WaitGroup counter
		go func(i int, ga artist.Artist) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			defer func() {
				if r := recover(); r != nil {
					select {
					case errs <- fmt.Errorf("error in goroutine: %v", r):
					default:
						// If the channel is already full, don't block
					}
				}
			}()
			result[i] = mapGormArtistToGqlArtist(&ga)
		}(i, gormArtist)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errs) // Close the channel to signal no more errors will be sent

	// Check for errors
	if err, ok := <-errs; ok {
		return nil, "", err
	}

	return result, nextCursor, nil
}

func (s *ArtistService) FindAllByCursor(ctx context.Context, cursor string, limit int) ([]*models.Artist, string, error) {
	// Implement your search logic
	artists, nextCursor, err := s.repo.FindAllByCursor(ctx, cursor, limit)

	if err != nil {
		return nil, "", err
	}

	var wg sync.WaitGroup
	result := make([]*models.Artist, len(artists))
	errs := make(chan error, 1) // Buffered channel to avoid blocking

	for i, gormArtist := range artists {
		wg.Add(1) // Increment the WaitGroup counter
		go func(i int, ga artist.Artist) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			defer func() {
				if r := recover(); r != nil {
					select {
					case errs <- fmt.Errorf("error in goroutine: %v", r):
					default:
						// If the channel is already full, don't block
					}
				}
			}()
			result[i] = mapGormArtistToGqlArtist(&ga)
		}(i, gormArtist)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errs) // Close the channel to signal no more errors will be sent

	// Check for errors
	if err, ok := <-errs; ok {
		return nil, "", err
	}

	return result, nextCursor, nil
}

func (s *ArtistService) Save(ctx context.Context, artist *models.Artist) (*models.Artist, error) {
	gormArtist := mapGqlArtistToGormArtist(artist)
	savedArtist, err := s.repo.Save(ctx, gormArtist)
	if err != nil {
		return nil, err
	}
	return mapGormArtistToGqlArtist(savedArtist), nil
}

func (s *ArtistService) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	success, err := s.repo.Delete(ctx, id)
	return success, err
}

func (s *ArtistService) Update(ctx context.Context, artist *models.Artist) (*models.Artist, error) {
	gormArtist := mapGqlArtistToGormArtist(artist)
	updatedArtist, err := s.repo.Update(ctx, gormArtist)
	if err != nil {
		return nil, err
	}
	return mapGormArtistToGqlArtist(updatedArtist), nil
}

func mapGqlArtistToGormArtist(gqlArtist *models.Artist) *artist.Artist {
	gqlArtistDto := &artist.Artist{
		ID:            gqlArtist.ID,
		Name:          gqlArtist.Name,
		Location:      *gqlArtist.Location,
		SCCity:        *gqlArtist.City,
		SCCountry:     *gqlArtist.Country,
		SCAvatarURL:   *gqlArtist.AvatarURL,
		SCFirstName:   *gqlArtist.FirstName,
		SCLastName:    *gqlArtist.LastName,
		SCFullName:    *gqlArtist.FullName,
		SCUsername:    gqlArtist.Username,
		SCDescription: *gqlArtist.Description,
		SCPromotedSet: *gqlArtist.SoundcloudPromotedSet,
	}

	for _, sm := range gqlArtist.SocialMediaLinks {
		gqlSocialMedia := artist.SocialMedia{
			ID:       sm.ID,
			Platform: sm.Platform,
			Link:     sm.Link,
		}
		gqlArtistDto.SocialMediaLinks = append(gqlArtistDto.SocialMediaLinks, gqlSocialMedia)
	}

	return gqlArtistDto
}

func mapGormArtistToGqlArtist(gormArtist *artist.Artist) *models.Artist {

	gqlArtist := &models.Artist{
		ID:                    gormArtist.ID,
		Name:                  gormArtist.Name,
		Location:              &gormArtist.Location,
		City:                  &gormArtist.SCCity,
		Country:               &gormArtist.SCCountry,
		AvatarURL:             &gormArtist.SCAvatarURL,
		FirstName:             &gormArtist.SCFirstName,
		LastName:              &gormArtist.SCLastName,
		FullName:              &gormArtist.SCFullName,
		Username:              gormArtist.SCUsername,
		Description:           &gormArtist.SCDescription,
		SoundcloudPromotedSet: &gormArtist.SCPromotedSet,
	}

	for _, sm := range gormArtist.SocialMediaLinks {
		gqlSocialMedia := models.SocialMedia{
			ID:       sm.ID, // Assuming UUID is used as the ID in GraphQL model
			Platform: sm.Platform,
			Link:     sm.Link,
		}
		gqlArtist.SocialMediaLinks = append(gqlArtist.SocialMediaLinks, &gqlSocialMedia)
	}

	return gqlArtist
}
