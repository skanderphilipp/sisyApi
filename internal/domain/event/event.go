package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
	"github.com/skanderphilipp/sisyApi/internal/domain/stage"
	"github.com/skanderphilipp/sisyApi/internal/domain/venue"
	"gorm.io/gorm"
)

type Event struct {
	ID        uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VenueID   uuid.UUID         `gorm:"type:uuid;not null" json:"venueID"`
	Venue     *venue.Venue      `gorm:"foreignKey:VenueID" json:"venue"`
	StartDate time.Time         `gorm:"not null" json:"startDate"`
	EndDate   time.Time         `gorm:"not null" json:"endDate"`
	Timetable []*TimetableEntry `gorm:"foreignKey:EventID" json:"timetable,omitempty"`
	//gorm additional fields
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type TimetableEntry struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EventID   uuid.UUID      `gorm:"type:uuid;foreignKey:EventID" json:"eventID"`
	StageID   uuid.UUID      `gorm:"type:uuid;foreignKey:StageID" json:"stageID"`
	Stage     *stage.Stage   `json:"stage,omitempty"`
	ArtistID  uuid.UUID      `gorm:"type:uuid;foreignKey:ArtistID" json:"artistID"`
	Artist    *artist.Artist `json:"artist,omitempty"`
	StartTime time.Time      `json:"startTime,omitempty"`
	EndTime   time.Time      `json:"endTime,omitempty"`
	//gorm additinonal fields
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// Additional fields like CreatedAt, UpdatedAt can be added.
}

func (e *Event) AddTimetableEntry(entry *TimetableEntry) error {
	// Validate entry
	if err := e.validateTimetableEntry(entry); err != nil {
		return err
	}
	e.Timetable = append(e.Timetable, entry)
	return nil
}
func (e *Event) validateTimetableEntry(entry *TimetableEntry) error {
	if entry.StartTime.Before(e.StartDate) || entry.EndTime.After(e.EndDate) {
		return fmt.Errorf("timetable entry times must be within the event start and end dates")
	}
	return nil
}

// Gorm Hooks

// Event BeforeCreate hook
func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

func (e *TimetableEntry) BeforeCreate(tx *gorm.DB) error {
	var event Event
	if err := tx.First(&event, "id = ?", e.EventID).Error; err != nil {
		return err
	}

	return event.validateTimetableEntry(e)
}

func (e *Event) BeforeUpdate(tx *gorm.DB) error {
	// Perform any necessary validation or updates
	// For example, you might want to validate related TimetableEntries
	// or update them based on changes to the Event.

	// Example validation (you can customize as needed):
	for _, entry := range e.Timetable {
		if err := e.validateTimetableEntry(entry); err != nil {
			return err
		}
	}

	return nil
}
