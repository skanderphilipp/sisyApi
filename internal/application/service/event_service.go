package service

import (
	"context"

	"github.com/blnto/blnto_service/internal/domain/event"
	"github.com/blnto/blnto_service/internal/domain/models"
	"github.com/blnto/blnto_service/internal/infrastructure/repository"
	"github.com/google/uuid"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func mapGormEventToGqlEvent(gormEvent *event.Event) *models.Event {
	if gormEvent == nil {
		return nil // or handle it according to your application's needs
	}

	gqlEvent := &models.Event{
		ID:        gormEvent.ID,
		StartDate: gormEvent.StartDate,
		EndDate:   gormEvent.EndDate,
	}

	if gormEvent.Venue != nil {
		gqlEvent.Venue = mapGormVenueToGqlVenue(gormEvent.Venue)
	}

	if gormEvent.Timetable != nil {
		gqlEvent.Timetable = mapGormTimetableEntriesToGql(gormEvent.Timetable)
	}

	return gqlEvent
}

func mapGormTimetableEntriesToGql(gormEntries []*event.TimetableEntry) []*models.TimetableEntry {
	var gqlEntries []*models.TimetableEntry
	for _, entry := range gormEntries {
		gqlEntry := &models.TimetableEntry{
			ID:        entry.ID,
			EventID:   entry.EventID,
			StageID:   entry.StageID,
			Stage:     mapGormStageToGqlStage(entry.Stage),
			ArtistID:  entry.ArtistID,
			Artist:    mapGormArtistToGqlArtist(entry.Artist),
			StartTime: &entry.StartTime,
			EndTime:   &entry.EndTime,
		}
		gqlEntries = append(gqlEntries, gqlEntry)
	}
	return gqlEntries
}

func mapGqlEventToGormEvent(gqlEvent *models.Event) *event.Event {
	gormEvent := &event.Event{
		ID:        gqlEvent.ID,
		StartDate: gqlEvent.StartDate,
		EndDate:   gqlEvent.EndDate,
		VenueID:   gqlEvent.Venue.ID,
		Timetable: mapGqlTimetableEntriesToGorm(gqlEvent.Timetable),
		// Venue:     mapGqlVenueToGormVenue(gqlEvent.Venue),
	}

	return gormEvent
}

func mapGqlTimetableEntriesToGorm(gqlEntries []*models.TimetableEntry) []*event.TimetableEntry {
	var gormEntries []*event.TimetableEntry
	for _, entry := range gqlEntries {
		gormEntry := &event.TimetableEntry{
			ID:        entry.ID,
			EventID:   entry.EventID,
			StageID:   entry.StageID,
			ArtistID:  entry.ArtistID,
			StartTime: *entry.StartTime,
			EndTime:   *entry.EndTime,
		}
		gormEntries = append(gormEntries, gormEntry)
	}
	return gormEntries
}

func (s *EventService) FindAllByVenueID(ctx context.Context, venueId uuid.UUID) ([]*models.Event, error) {
	gormEvents, err := s.repo.FindAllByVenueID(ctx, venueId)
	if err != nil {
		return nil, err
	}

	var gqlEvents []*models.Event
	for _, eventData := range gormEvents {
		gqlEvents = append(gqlEvents, mapGormEventToGqlEvent(eventData))
	}

	return gqlEvents, nil
}
func (s *EventService) FindUpcomingByVenueID(ctx context.Context, venueID uuid.UUID) ([]*models.Event, error) {
	gormEvents, err := s.repo.FindUpcomingByVenueID(ctx, venueID)
	if err != nil {
		return nil, err
	}

	var gqlEvents []*models.Event
	for _, eventData := range gormEvents {
		gqlEvents = append(gqlEvents, mapGormEventToGqlEvent(eventData))
	}

	return gqlEvents, nil
}

func (s *EventService) FindPastEventsByVenueID(ctx context.Context, venueID uuid.UUID) ([]*models.Event, error) {
	events, err := s.repo.FindPastEventsByVenueID(ctx, venueID)
	if err != nil {
		return nil, err
	}
	var result []*models.Event
	for _, eventData := range events {
		result = append(result, mapGormEventToGqlEvent(eventData))
	}

	return result, err
}

func (s *EventService) FindAllUpcoming(ctx context.Context, cursor string, limit int) ([]*models.Event, string, error) {
	events, nextCursor, err := s.repo.FindAllUpcoming(ctx, cursor, limit)
	if err != nil {
		return nil, "", err
	}
	var result []*models.Event
	for _, eventData := range events {
		result = append(result, mapGormEventToGqlEvent(eventData))
	}

	return result, nextCursor, nil
}

func (s *EventService) FindToday(ctx context.Context) ([]*models.Event, error) {
	events, err := s.repo.FindToday(ctx)
	if err != nil {
		return nil, err
	}
	var result []*models.Event
	for _, eventData := range events {
		result = append(result, mapGormEventToGqlEvent(eventData))
	}
	return result, nil
}

func (s *EventService) FindTomorrow(ctx context.Context) ([]*models.Event, error) {
	events, err := s.repo.FindTomorrow(ctx)
	if err != nil {
		return nil, err
	}
	var result []*models.Event
	for _, eventData := range events {
		result = append(result, mapGormEventToGqlEvent(eventData))
	}
	return result, nil
}

func (s *EventService) FindCurrent(ctx context.Context) ([]*models.Event, error) {
	events, err := s.repo.FindCurrent(ctx)
	if err != nil {
		return nil, err
	}
	var result []*models.Event
	for _, eventData := range events {
		result = append(result, mapGormEventToGqlEvent(eventData))
	}
	return result, nil
}

func (s *EventService) Save(ctx context.Context, gqlEvent *models.Event) (*models.Event, error) {
	gormEvent := mapGqlEventToGormEvent(gqlEvent)
	savedEvent, err := s.repo.Save(ctx, gormEvent)
	if err != nil {
		return nil, err
	}
	return mapGormEventToGqlEvent(savedEvent), nil
}

func (s *EventService) Update(ctx context.Context, gqlEvent *models.Event) (*models.Event, error) {
	gormEvent := mapGqlEventToGormEvent(gqlEvent)
	updatedEvent, err := s.repo.Update(ctx, gormEvent)
	if err != nil {
		return nil, err
	}
	return mapGormEventToGqlEvent(updatedEvent), nil
}

func (s *EventService) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	success, err := s.repo.Delete(ctx, id)
	return success, err
}
