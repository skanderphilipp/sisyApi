package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/blnto/blnto_service/internal/domain/event"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (repo *EventRepository) FindUpcomingByVenueID(ctx context.Context, venueID uuid.UUID) ([]*event.Event, error) {
	var events []*event.Event
	err := repo.db.WithContext(ctx).Where("venue_id = ? AND start_date > ?", venueID, time.Now()).
		Preload("Venue.Stages").
		Preload("Timetable.Stage").
		Preload("Timetable.Artist").
		Find(&events).Error
	return events, err
}

func (repo *EventRepository) FindPastEventsByVenueID(ctx context.Context, venueID uuid.UUID) ([]*event.Event, error) {
	var events []*event.Event
	err := repo.db.WithContext(ctx).Where("venue_id = ? AND end_date < ?", venueID, time.Now()).
		Preload("Venue.Stages").
		Preload("Timetable.Stage").
		Preload("Timetable.Artist").
		Find(&events).Error
	return events, err
}
func (repo *EventRepository) FindAllByVenueID(ctx context.Context, venueId uuid.UUID) ([]*event.Event, error) {
	var events []*event.Event
	err := repo.db.WithContext(ctx).Where("venue_id = ?", venueId).
		Preload("Venue.Stages").
		Preload("Timetable.Stage").
		Preload("Timetable.Artist").
		Find(&events).Error
	return events, err
}

func (repo *EventRepository) FindAllUpcoming(ctx context.Context, cursort string, limit int) ([]*event.Event, string, error) {
	var events []*event.Event
	var nextCursor string

	err := repo.db.WithContext(ctx).Where("start_date > ?", time.Now()).
		Preload("Venue.Stages").
		Preload("Timetable.Stage").
		Preload("Timetable.Artist").
		Find(&events).Error

	if err != nil {
		return nil, "", err
	}

	// Set next cursor
	if len(events) > 0 {
		nextCursor = events[len(events)-1].ID.String() // Assuming ID is a UUID
	}

	return events, nextCursor, nil
}

func (repo *EventRepository) FindToday(ctx context.Context) ([]*event.Event, error) {
	today := time.Now()
	var events []*event.Event
	err := repo.db.WithContext(ctx).Where("DATE(start_date) = DATE(?)", today).
		Preload("Venue.Stages").
		Preload("Timetable.Stage").
		Preload("Timetable.Artist").
		Find(&events).Error
	return events, err
}

func (repo *EventRepository) FindTomorrow(ctx context.Context) ([]*event.Event, error) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	var events []*event.Event
	err := repo.db.WithContext(ctx).Where("DATE(start_date) = DATE(?)", tomorrow).Preload("Venue.Stages").
		Preload("Timetable.Stage").
		Preload("Timetable.Artist").
		Find(&events).Error
	return events, err
}

func (repo *EventRepository) FindCurrent(ctx context.Context) ([]*event.Event, error) {
	now := time.Now()
	var events []*event.Event
	err := repo.db.WithContext(ctx).Where("start_date <= ? AND end_date >= ?", now, now).
		Preload("Venue.Stages").
		Preload("Timetable.Stage").
		Preload("Timetable.Artist").
		Find(&events).Error
	return events, err
}

func (repo *EventRepository) Save(ctx context.Context, event *event.Event) (*event.Event, error) {
	result := repo.db.WithContext(ctx).Save(event)

	if result.Error != nil {
		return nil, fmt.Errorf("error saving artist: %v", result.Error)
	}

	return event, nil
}

func (repo *EventRepository) Update(ctx context.Context, eventData *event.Event) (*event.Event, error) {
	if err := repo.db.WithContext(ctx).Save(eventData).Error; err != nil {
		return nil, err
	}
	return eventData, nil
}

func (repo *EventRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	// Delete the artist with the given ID from the database
	result := repo.db.WithContext(ctx).Delete(&event.Event{}, id)

	if result.Error != nil {
		return false, fmt.Errorf("error deleting artist: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("artist not found")
	}

	return true, nil
}
