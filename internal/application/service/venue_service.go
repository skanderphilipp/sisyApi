package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/blnto/blnto_service/internal/domain/models"
	"github.com/blnto/blnto_service/internal/domain/stage"
	"github.com/blnto/blnto_service/internal/domain/venue"
	"github.com/blnto/blnto_service/internal/infrastructure/repository"
	"github.com/google/uuid"
)

type VenueService struct {
	repo *repository.VenueRepository
}

func NewVenueService(repo *repository.VenueRepository) *VenueService {
	return &VenueService{repo: repo}
}

func (s *VenueService) FindAllByCursor(ctx context.Context, cursor string, limit int) ([]*models.Venue, string, error) {
	// Implement your search logic
	venues, nextCursor, err := s.repo.FindAllByCursor(ctx, cursor, limit)

	if err != nil {
		return nil, "", err
	}

	var wg sync.WaitGroup
	result := make([]*models.Venue, len(venues))
	errs := make(chan error, 1) // Buffered channel to avoid blocking

	for i, gormVenue := range venues {
		wg.Add(1) // Increment the WaitGroup counter
		go func(i int, ga venue.Venue) {
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
			result[i] = mapGormVenueToGqlVenue(&ga)
		}(i, gormVenue)
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

func (s *VenueService) Save(ctx context.Context, gqlVenue *models.Venue) (*models.Venue, error) {
	gormVenue := mapGqlVenueToGormVenue(gqlVenue)
	savedVenue, err := s.repo.Save(ctx, gormVenue)
	if err != nil {
		return nil, err
	}
	return mapGormVenueToGqlVenue(savedVenue), nil
}

func (s *VenueService) Update(ctx context.Context, gqlVenue *models.Venue) (*models.Venue, error) {
	gormVenue := mapGqlVenueToGormVenue(gqlVenue)
	updatedVenue, err := s.repo.Update(ctx, gormVenue)
	if err != nil {
		return nil, err
	}
	return mapGormVenueToGqlVenue(updatedVenue), nil
}

func (s *VenueService) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	success, err := s.repo.Delete(ctx, id)
	return success, err
}

func mapGormVenueToGqlVenue(gormVenue *venue.Venue) *models.Venue {
	gqlVenue := &models.Venue{
		ID:          gormVenue.ID,
		Name:        gormVenue.Name,
		Description: gormVenue.Description,
	}
	// Check if there are stages to map
	if len(gormVenue.Stages) > 0 {
		for _, stageData := range gormVenue.Stages {
			// Assuming that models.Stage and venue.Stage have similar fields
			gqlStage := &models.Stage{
				ID:   stageData.ID,
				Name: stageData.StageName,
			}
			gqlVenue.Stages = append(gqlVenue.Stages, gqlStage)
		}
	}
	return gqlVenue
}

func mapGqlVenueToGormVenue(gqlVenue *models.Venue) *venue.Venue {
	gormVenue := &venue.Venue{
		ID:          gqlVenue.ID,
		Name:        gqlVenue.Name,
		Description: gqlVenue.Description,
	}

	// Check if there are stages to map
	if len(gqlVenue.Stages) > 0 {
		for _, stageData := range gqlVenue.Stages {
			// Assuming that models.Stage and venue.Stage have similar fields
			gormStage := &stage.Stage{
				StageName: stageData.Name,
			}
			gormVenue.Stages = append(gormVenue.Stages, gormStage)
		}
	}

	return gormVenue
}
