package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
)

// CreateVenue is the resolver for the createVenue field.
func (r *mutationResolver) CreateVenue(ctx context.Context, input models.CreateVenueInput) (*models.Venue, error) {
	panic(fmt.Errorf("not implemented: CreateVenue - createVenue"))
}

// ListVenues is the resolver for the listVenues field.
func (r *queryResolver) ListVenues(ctx context.Context, first *int, after *string) (*models.VenueConnection, error) {
	// Set default values if nil
	var limit int
	if first == nil {
		limit = 10 // Default limit value
	} else {
		limit = *first
	}

	var cursor string
	if after != nil {
		cursor = *after
	}

	// Assuming you have an artist repository instance (artistRepo)
	venues, nextCursor, err := r.venueRepo.FindAllByCursor(ctx, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching artists: %v", err)
	}

	// Map artists to GraphQL edges

	edges := make([]*models.VenueEdge, len(venues))

	for i, venue := range venues {
		cursorStr := venue.ID.String() // Convert UUID to string
		edges[i] = &models.VenueEdge{
			Node:   &venue,
			Cursor: cursorStr, // Assuming ID is a UUID
		}
	}

	// Construct ArtistConnection
	hasNextPage := len(edges) == limit
	venueConnection := &models.VenueConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			EndCursor:   &nextCursor,
			HasNextPage: &hasNextPage,
		},
	}

	return venueConnection, nil
}

// AllVenues is the resolver for the allVenues field.
func (r *queryResolver) AllVenues(ctx context.Context) ([]*models.Venue, error) {
	panic(fmt.Errorf("not implemented: AllVenues - allVenues"))
}

// Venue is the resolver for the venue field.
func (r *queryResolver) Venue(ctx context.Context, id uuid.UUID) (*models.Venue, error) {
	panic(fmt.Errorf("not implemented: Venue - venue"))
}
