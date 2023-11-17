package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (repo *EventRepository) FindUpcomingByVenueID(ctx context.Context, venueID uuid.UUID) ([]models.Event, error) {
	var events []models.Event
	err := repo.db.WithContext(ctx).Where("venue_id = ? AND start_date > ?", venueID, time.Now()).Find(&events).Error
	return events, err
}

func (repo *EventRepository) FindPastEventsByVenueID(ctx context.Context, venueID uuid.UUID) ([]models.Event, error) {
	var events []models.Event
	err := repo.db.WithContext(ctx).Where("venue_id = ? AND end_date < ?", venueID, time.Now()).Find(&events).Error
	return events, err
}

func (repo *EventRepository) FindAllUpcoming(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	err := repo.db.WithContext(ctx).Where("start_date > ?", time.Now()).Find(&events).Error
	return events, err
}

func (repo *EventRepository) FindToday(ctx context.Context) ([]models.Event, error) {
	today := time.Now()
	var events []models.Event
	err := repo.db.WithContext(ctx).Where("DATE(start_date) = DATE(?)", today).Find(&events).Error
	return events, err
}

func (repo *EventRepository) FindTomorrow(ctx context.Context) ([]models.Event, error) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	var events []models.Event
	err := repo.db.WithContext(ctx).Where("DATE(start_date) = DATE(?)", tomorrow).Find(&events).Error
	return events, err
}

func (repo *EventRepository) FindCurrent(ctx context.Context) ([]models.Event, error) {
	now := time.Now()
	var events []models.Event
	err := repo.db.WithContext(ctx).Where("start_date <= ? AND end_date >= ?", now, now).Find(&events).Error
	return events, err
}

func (repo *EventRepository) Save(ctx context.Context, event *models.Event) error {
	return repo.db.WithContext(ctx).Save(event).Error
}

func (repo *EventRepository) Update(ctx context.Context, event *models.Event) error {
	return repo.db.WithContext(ctx).Model(&models.Event{}).Where("id = ?", event.ID).Updates(event).Error
}

func (repo *EventRepository) Delete(ctx context.Context, eventID uuid.UUID) error {
	return repo.db.WithContext(ctx).Delete(&models.Event{}, "id = ?", eventID).Error
}
