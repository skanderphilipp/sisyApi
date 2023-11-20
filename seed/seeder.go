package seed

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
	"github.com/skanderphilipp/sisyApi/internal/domain/event"
	"github.com/skanderphilipp/sisyApi/internal/domain/stage"
	"github.com/skanderphilipp/sisyApi/internal/domain/venue"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		}
	}()

	var venues []venue.Venue
	var stages []stage.Stage
	var artists []artist.Artist
	var err error

	if venues, err = seedVenues(tx, 2); err != nil {
		return err
	}

	if stages, err = seedStages(tx, venues, 20); err != nil {
		return err
	}

	if err = seedEventsAndTimetableEntries(tx, venues, stages, artists, 20, 20); err != nil {
		return err
	}

	return tx.Commit().Error
}

func seedVenues(tx *gorm.DB, count int) ([]venue.Venue, error) {
	venues := make([]venue.Venue, count)
	for i := range venues {
		sentence := gofakeit.Sentence(10)
		venues[i].ID = uuid.New()
		venues[i].Name = gofakeit.Name()
		venues[i].Description = &sentence
	}

	if err := tx.Create(&venues).Error; err != nil {
		return nil, err
	}

	return venues, nil
}

func seedStages(tx *gorm.DB, venues []venue.Venue, count int) ([]stage.Stage, error) {
	stages := make([]stage.Stage, count)
	for i := range stages {
		stages[i].ID = uuid.New()
		stages[i].StageName = gofakeit.Word()
		stages[i].VenueID = venues[i%len(venues)].ID
	}

	if err := tx.Create(&stages).Error; err != nil {
		return nil, err
	}

	return stages, nil
}

func seedEventsAndTimetableEntries(tx *gorm.DB, venues []venue.Venue, stages []stage.Stage, artists []artist.Artist, eventCount, timetableEntryCount int) error {
	for i := 0; i < eventCount; i++ {
		// Generate random start and end dates for the event
		startDate, endDate := generateRandomDates()
		eventDto := event.Event{
			ID:        uuid.New(),
			VenueID:   venues[i%len(venues)].ID,
			StartDate: startDate,
			EndDate:   endDate,
		}

		for j := 0; j < timetableEntryCount; j++ {
			entryStartTime, entryEndTime := generateRandomTimeEntries(startDate, endDate)

			timetableEntry := &event.TimetableEntry{ // Create a pointer
				EventID:   eventDto.ID,
				StageID:   stages[j%len(stages)].ID,
				ArtistID:  artists[j%len(artists)].ID,
				StartTime: entryStartTime,
				EndTime:   entryEndTime,
			}

			eventDto.Timetable = append(eventDto.Timetable, timetableEntry)
		}

		// Create the event with its timetable entries after the inner loop
		if err := tx.Create(&eventDto).Error; err != nil {
			return err
		}
	}
	return nil
}

// generateRandomDates returns random start and end dates
func generateRandomDates() (time.Time, time.Time) {
	now := time.Now()
	randomDays := gofakeit.Number(-100, 100) // Random number of days up to +/- 1 year
	startDate := now.AddDate(0, 0, randomDays)
	endDate := startDate.AddDate(0, 0, gofakeit.Number(1, 4)) // Event lasts between 1 to 30 days
	return startDate, endDate
}

// generateRandomEntryDates returns random start and end times for a timetable entry within the given date range
func generateRandomTimeEntries(startDate, endDate time.Time) (time.Time, time.Time) {
	entryStart := gofakeit.DateRange(startDate, endDate)
	entryEnd := entryStart.Add(time.Duration(gofakeit.Number(1, 4)) * time.Hour) // Lasts 1 to 4 hours
	if entryEnd.After(endDate) {
		entryEnd = endDate
	}
	return entryStart, entryEnd
}
