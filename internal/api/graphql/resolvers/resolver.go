package resolvers

import (
	"github.com/skanderphilipp/sisyApi/internal/application/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	artistService *service.ArtistService
	eventService  *service.EventService
	stageService  *service.StageService
	venueService  *service.VenueService
}

func NewResolver(artistService *service.ArtistService, eventService *service.EventService, stageService *service.StageService, venueService *service.VenueService) *Resolver {
	return &Resolver{artistService: artistService, eventService: eventService, stageService: stageService, venueService: venueService}
}
