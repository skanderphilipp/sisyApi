package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
	graphql1 "github.com/skanderphilipp/sisyApi/internal/infrastructure/graphql"
	"github.com/skanderphilipp/sisyApi/internal/utils"
)

// CreateArtist is the resolver for the createArtist field.
func (r *mutationResolver) CreateArtist(ctx context.Context, input models.CreateArtistInput) (*models.Artist, error) {
	newArtist := models.Artist{
		Name:                  input.Name,
		Location:              input.Location,
		SoundcloudPromotedSet: input.SoundcloudPromotedSet,
	}

	// Call the repository function to create the artist
	createdArtist, err := r.artistService.Save(ctx, &newArtist)
	if err != nil {
		return nil, err
	}

	return createdArtist, nil
}

// UpdateArtist is the resolver for the updateArtist field.
func (r *mutationResolver) UpdateArtist(ctx context.Context, input models.UpdateArtistInput) (*models.Artist, error) {
	panic(fmt.Errorf("not implemented: UpdateArtist - updateArtist"))
}

// DeleteArtist is the resolver for the deleteArtist field.
func (r *mutationResolver) DeleteArtist(ctx context.Context, input models.DeleteArtistInput) (bool, error) {
	// Delete the artist by ID
	result, err := r.artistService.Delete(ctx, input.ID)
	if err != nil {
		return false, err
	}

	return result, nil
}

// GetArtist is the resolver for the getArtist field.
func (r *queryResolver) GetArtist(ctx context.Context, id uuid.UUID) (*models.Artist, error) {
	artist, err := r.artistService.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching artist: %v", err)
	}
	return artist, err
}

// SearchArtists is the resolver for the searchArtists field.
func (r *queryResolver) SearchArtists(ctx context.Context, criteria models.ArtistSearchInput) (*models.ArtistConnection, error) {
	// Set default values if nil
	if criteria.First == nil {
		defaultValue := 10             // Default value
		criteria.First = &defaultValue // Assign the address of the default value
	}

	// Assuming you have an artist repository instance (artistRepo)
	artists, nextCursor, err := r.artistService.Search(ctx, &criteria)
	if err != nil {
		return nil, fmt.Errorf("error fetching artists: %v", err)
	}

	// Map artists to GraphQL edges

	edges := make([]*models.ArtistEdge, len(artists))

	for i, artist := range artists {
		cursorStr := artist.ID.String() // Convert UUID to string
		edges[i] = &models.ArtistEdge{
			Node:   artist,
			Cursor: &cursorStr, // Assuming ID is a UUID
		}
	}

	// Construct ArtistConnection
	hasNextPage := len(edges) == *criteria.First
	artistConnection := &models.ArtistConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			EndCursor:   &nextCursor,
			HasNextPage: &hasNextPage,
		},
	}

	return artistConnection, nil
}

// GetArtistByName is the resolver for the getArtistByName field.
func (r *queryResolver) GetArtistByName(ctx context.Context, name string) (*models.Artist, error) {
	artist, err := r.artistService.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error fetching artist: %v", err)
	}

	return artist, nil
}

// ListArtists is the resolver for the listArtists field.
func (r *queryResolver) ListArtists(ctx context.Context, first *int, after *string) (*models.ArtistConnection, error) {
	cursorFunc := func(artist *models.Artist) string {
		return artist.ID.String() // Assuming ID is a UUID
	}

	artists, nextCursor, limit, err := utils.FetchItemsList[models.Artist](ctx, first, after, r.artistService.FindAllByCursor)

	if err != nil {
		return nil, fmt.Errorf("error fetching artists: %v", err)
	}
	// Map items to GraphQL edges
	edges := make([]*models.ArtistEdge, len(artists))

	for i, item := range artists {
		cursorStr := cursorFunc(item)
		edges[i] = &models.ArtistEdge{
			Node:   item,
			Cursor: &cursorStr,
		}
	}

	// Construct ArtistConnection
	hasNextPage := len(edges) == limit
	artistConnection := &models.ArtistConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			EndCursor:   &nextCursor,
			HasNextPage: &hasNextPage,
		},
	}

	return artistConnection, nil
}

// Mutation returns graphql1.MutationResolver implementation.
func (r *Resolver) Mutation() graphql1.MutationResolver { return &mutationResolver{r} }

// Query returns graphql1.QueryResolver implementation.
func (r *Resolver) Query() graphql1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
