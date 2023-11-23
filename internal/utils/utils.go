package utils

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/blnto/blnto_service/internal/domain/models"
)

type Pagination struct {
	Limit        int    `json:"limit"`
	CurrentPage  int    `json:"currentPage"`
	NextPage     string `json:"nextPage"`
	PreviousPage string `json:"previousPage"`
	SortField    string `json:"sortField"`
	SortOrder    string `json:"sortOrder"`
	TotalRecords int    `json:"totalRecords"`
}

func FetchItemsList[T any](ctx context.Context, first *int, after *string, fetchFunc func(context.Context, string, int) ([]*T, string, error)) ([]*T, string, int, error) {
	// Set default values if nil
	limit := 10 // Default limit
	if first != nil {
		limit = *first
	}

	cursor := ""
	if after != nil {
		cursor = *after
	}

	// Fetch items using the provided fetch function
	items, nextCursor, err := fetchFunc(ctx, cursor, limit)

	if err != nil {
		return nil, "", 10, fmt.Errorf("error fetching items: %v", err)
	}

	return items, nextCursor, limit, nil
}

func BuildEventConnection(events []*models.Event, limit int, nextCursor string, cursorFunc func(*models.Event) string) *models.EventConnection {
	edges := make([]*models.EventEdge, len(events))

	for i, item := range events {
		cursorStr := cursorFunc(item)
		edges[i] = &models.EventEdge{
			Node:   item,
			Cursor: &cursorStr,
		}
	}

	hasNextPage := len(edges) == limit
	return &models.EventConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			EndCursor:   &nextCursor,
			HasNextPage: &hasNextPage,
		},
	}
}

func IsValidURL(toTest string) bool {
	// This is a simple check. For more complex validation, you might want to use regex.
	u, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	// Ensure the scheme is http or https
	scheme := strings.ToLower(u.Scheme)
	return scheme == "http" || scheme == "https"
}
