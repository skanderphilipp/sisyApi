package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/blnto/blnto_service/internal/domain/artist"
	"github.com/blnto/blnto_service/internal/domain/models"
	"github.com/blnto/blnto_service/internal/infrastructure/repository"
	"github.com/google/uuid"
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
	artistModel, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return mapGormArtistToGqlArtist(artistModel), nil
}

func (s *ArtistService) FindByID(ctx context.Context, id uuid.UUID) (*models.Artist, error) {
	artistModel, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapGormArtistToGqlArtist(artistModel), nil
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
	gormArtist := s.createGormArtistFromGqlArtist(artist)

	// Validate and check uniqueness of permalink
	if err := s.validateAndCheckPermalink(ctx, artist); err != nil {
		return nil, err
	}
	// Set permalink
	if err := gormArtist.SetPermalink(artist.SoundcloudPermalink); err != nil {
		return nil, fmt.Errorf("failed to set permalink: %w", err)
	}

	// Validate social media links
	for _, link := range gormArtist.SocialMediaLinks {
		err := link.ValidateLinkFormat()
		if err != nil {
			return nil, err
		}
	}

	if err := s.validateSocialMediaPlatforms(gormArtist.SocialMediaLinks); err != nil {
		return nil, err
	}

	// Save the artist
	savedArtist, err := s.repo.Save(ctx, gormArtist)
	if err != nil {
		return nil, fmt.Errorf("failed to save artist: %w", err)
	}

	return mapGormArtistToGqlArtist(savedArtist), nil
}

func (s *ArtistService) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	success, err := s.repo.Delete(ctx, id)
	return success, err
}

func (s *ArtistService) Update(ctx context.Context, artist *models.Artist) (*models.Artist, error) {
	// Validate and check uniqueness of permalink
	if err := s.validateAndCheckPermalink(ctx, artist); err != nil {
		return nil, err
	}

	// Fetch the existing artist
	existingArtist, err := s.repo.FindByID(ctx, artist.ID)
	if err != nil {
		return nil, err
	}

	// Set permalink using the SetPermalink method
	if err := existingArtist.SetPermalink(artist.SoundcloudPermalink); err != nil {
		return nil, fmt.Errorf("failed to set permalink: %w", err)
	}

	// Update fields from artist to existingArtist...
	s.updateGormArtistFromGqlArtist(existingArtist, artist)

	// Handle social media links update
	if err := s.handleSocialMediaLinksUpdate(ctx, existingArtist, artist); err != nil {
		return nil, fmt.Errorf("failed to update social media links: %w", err)
	}

	// Save the updated artist
	updatedArtist, err := s.repo.Update(ctx, existingArtist)
	if err != nil {
		return nil, err
	}

	return mapGormArtistToGqlArtist(updatedArtist), nil
}

func (s *ArtistService) handleSocialMediaLinksUpdate(ctx context.Context, existingArtist *artist.Artist, gqlArtist *models.Artist) error {
	// Fetch current social media links from the database for existingArtist
	currentLinks, err := s.repo.GetSocialMediaLinksByArtistID(ctx, existingArtist.ID)
	if err != nil {
		return err
	}

	// Map to easily check existing links
	currentLinksMap := make(map[uuid.UUID]*artist.SocialMediaLink)
	for _, link := range currentLinks {
		currentLinksMap[link.ID] = link
	}

	// Loop through the provided artist's social media links
	for _, link := range gqlArtist.SocialMediaLinks {
		if link.ID == uuid.Nil {
			// New link, create it
			err := s.createSocialMediaLink(ctx, existingArtist.ID, link)
			if err != nil {
				return err
			}
		} else {
			// Existing link, update it
			err := s.updateSocialMediaLink(ctx, link)
			if err != nil {
				return err
			}
			// Remove the link from the map as it's still present
			delete(currentLinksMap, link.ID)
		}
	}

	// Delete any links that are no longer present
	for id := range currentLinksMap {
		err := s.repo.DeleteSocialMediaLink(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *ArtistService) createSocialMediaLink(ctx context.Context, artistID uuid.UUID, link *models.SocialMedia) error {
	// Implementation for creating a new social media link
	newLink := &artist.SocialMediaLink{
		ArtistID: artistID,
		Platform: artist.SocialMediaPlatform(link.Platform),
		Link:     link.Link,
	}
	return s.repo.CreateSocialMediaLink(ctx, *newLink)
}

func (s *ArtistService) updateSocialMediaLink(ctx context.Context, link *models.SocialMedia) error {
	// Implementation for updating an existing social media link
	existingLink := &artist.SocialMediaLink{
		ID:       link.ID,
		ArtistID: link.ArtistID,
		Platform: artist.SocialMediaPlatform(link.Platform),
		Link:     link.Link,
	}
	return s.repo.UpdateSocialMediaLink(ctx, *existingLink)
}

func (s *ArtistService) validateAndCheckPermalink(ctx context.Context, artist *models.Artist) error {
	if artist.SoundcloudPermalink != nil {
		if *artist.SoundcloudPermalink == "" {
			return errors.New("permalink cannot be an empty string")
		}
		exists, err := s.repo.PermalinkExistsExcludingArtist(ctx, *artist.SoundcloudPermalink, artist.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("permalink must be unique")
		}
	}
	return nil
}

func (s *ArtistService) validateSocialMediaPlatforms(links []artist.SocialMediaLink) error {
	allowedPlatforms := map[artist.SocialMediaPlatform]bool{
		artist.Twitter:    true,
		artist.Facebook:   true,
		artist.Instagram:  true,
		artist.YouTube:    true,
		artist.Soundcloud: true,
	}

	for _, link := range links {
		if _, ok := allowedPlatforms[link.Platform]; !ok {
			return fmt.Errorf("platform '%s' is not allowed", link.Platform)
		}
	}
	return nil
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
		gqlSocialMedia := artist.SocialMediaLink{
			ID:       sm.ID,
			Platform: artist.SocialMediaPlatform(sm.Platform),
			Link:     sm.Link,
		}
		gqlArtistDto.SocialMediaLinks = append(gqlArtistDto.SocialMediaLinks, gqlSocialMedia)
	}

	return gqlArtistDto
}

func (s *ArtistService) createGormArtistFromGqlArtist(gqlArtist *models.Artist) *artist.Artist {
	gormArtist := &artist.Artist{
		Name: gqlArtist.Name,
	}

	if gqlArtist.Location != nil {
		gormArtist.Location = *gqlArtist.Location
	}

	if gqlArtist.SoundcloudPromotedSet != nil {
		gormArtist.SCPromotedSet = *gqlArtist.SoundcloudPromotedSet
	}

	for _, sm := range gqlArtist.SocialMediaLinks {
		gormSocialMedia := artist.SocialMediaLink{
			Platform: artist.SocialMediaPlatform(sm.Platform),
			Link:     sm.Link,
		}
		gormArtist.SocialMediaLinks = append(gormArtist.SocialMediaLinks, gormSocialMedia)
	}

	return gormArtist
}

func (s *ArtistService) updateGormArtistFromGqlArtist(gormArtist *artist.Artist, gqlArtist *models.Artist) {
	gormArtist.Name = gqlArtist.Name

	if gqlArtist.Location != nil {
		gormArtist.Location = *gqlArtist.Location
	}

	if gqlArtist.SoundcloudPromotedSet != nil {
		gormArtist.SCPromotedSet = *gqlArtist.SoundcloudPromotedSet
	}
}

func (s *ArtistService) FindFeatured(ctx context.Context) ([]*models.Artist, error) {
	artists, err := s.repo.FindFeatured(ctx)

	if err != nil {
		return nil, err
	}

	var gormArtists []*models.Artist

	for _, artistDto := range artists {
		gormArtists = append(gormArtists, mapGormArtistToGqlArtist(&artistDto))
	}

	return gormArtists, nil
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
		SoundcloudID:          gormArtist.SCID,
		SoundcloudPromotedSet: &gormArtist.SCPromotedSet,
		SoundcloudPermalink:   gormArtist.SCPermalink,
	}

	for _, sm := range gormArtist.SocialMediaLinks {
		gqlSocialMedia := models.SocialMedia{
			ID:       sm.ID, // Assuming UUID is used as the ID in GraphQL model
			Platform: models.SocialMediaPlatform(sm.Platform),
			Link:     sm.Link,
		}
		gqlArtist.SocialMediaLinks = append(gqlArtist.SocialMediaLinks, &gqlSocialMedia)
	}

	return gqlArtist
}
