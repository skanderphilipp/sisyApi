package resolvers

import (
	"github.com/skanderphilipp/sisyApi/internal/application/service"
	"github.com/skanderphilipp/sisyApi/internal/infrastructure/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	artistService *service.ArtistService
	artistRepo    *repository.ArtistRepository
	venueRepo     *repository.VenueRepository
	eventRepo     *repository.EventRepository
}

func NewResolver(artistService *service.ArtistService, artistRepo *repository.ArtistRepository, venueRepo *repository.VenueRepository, eventRepo *repository.EventRepository) *Resolver {
	return &Resolver{artistRepo: artistRepo, artistService: artistService, venueRepo: venueRepo, eventRepo: eventRepo}
}
