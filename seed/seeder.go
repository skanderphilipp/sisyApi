package seed

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	artistModel "github.com/skanderphilipp/sisyApi/internal/domain/artist"
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
	"gorm.io/gorm"
)

func runSeedersIfNeeded(db *gorm.DB) {
	var isSeeded bool
	// Assume 'settings' is a table and 'seeded' is a column that stores a boolean
	db.Raw("SELECT seeded FROM settings").Scan(&isSeeded)

	if !isSeeded {
		// Run your seeders here
		// ...
		SeedArtistData(db)
		SeedStages(db)
		SeedTimetableData(db)
		// Update the flag
		db.Exec("UPDATE settings SET seeded = true")
	}
}

func SeedArtistData(db *gorm.DB) {
	// Initialize the random seed for gofakeit
	gofakeit.Seed(0)

	for i := 0; i < 150; i++ {
		// Generate fake data for each Artist
		artist := artistModel.Artist{
			UUID:              uuid.New(),
			Name:              gofakeit.Name(),
			Location:          gofakeit.City(),
			SoundcloudSetLink: gofakeit.URL(),
		}

		result := db.Create(&artist)
		if result.Error != nil {
			log.Fatalf("Failed to seed artist data: %v", result.Error)
		}

		// Generate fake data for associated SocialMedia
		// Randomly decide the number of social media links (for variety)
		numLinks := rand.Intn(5) + 1 // Ensure at least 1 link
		for j := 0; j < numLinks; j++ {
			socialMedia := artistModel.SocialMedia{
				UUID:                uuid.New(),
				ArtistID:            artist.UUID,
				SocialMediaPlatform: randomSocialMediaPlatform(),
				Link:                gofakeit.URL(),
			}

			result = db.Create(&socialMedia)
			if result.Error != nil {
				log.Fatalf("Failed to seed social media data: %v", result.Error)
			}
		}
	}
}

// randomSocialMediaPlatform returns a random social media platform
func randomSocialMediaPlatform() string {
	platforms := []string{
		string(artistModel.Twitter),
		string(artistModel.Facebook),
		string(artistModel.Instagram),
		string(artistModel.YouTube),
		string(artistModel.Soundcloud),
	}
	rand.Seed(time.Now().UnixNano())
	return platforms[rand.Intn(len(platforms))]
}

func SeedStages(db *gorm.DB) {
	// Create a random seed for generating UUIDs
	rand.Seed(time.Now().UnixNano())

	// Define the stage names
	stageNames := []string{"Hammahalle", "Wintergarten", "Dampfer", "Tunnel", "Strand"}

	// Iterate through the stage names and create records
	for _, name := range stageNames {
		stage := models.Stage{
			ID:        uuid.New(),
			StageName: name,
		}

		result := db.Create(&stage)

		if result.Error != nil {
			fmt.Printf("Error seeding stage %s: %v\n", name, result.Error)
		} else {
			fmt.Printf("Seeded stage: %s\n", name)
		}
	}
}

func SeedTimetableData(db *gorm.DB) {
	// Initialize the random seed for gofakeit
	gofakeit.Seed(0)

	// Create maps to store generated UUIDs for Artists and Stages
	artistUUIDs := make(map[string]uuid.UUID)

	for i := 0; i < 150; i++ {
		// Generate fake data for each Artist
		artist := artistModel.Artist{
			UUID:              uuid.New(),
			Name:              gofakeit.Name(),
			Location:          gofakeit.City(),
			SoundcloudSetLink: gofakeit.URL(),
		}

		result := db.Create(&artist)
		if result.Error != nil {
			log.Fatalf("Failed to seed artist data: %v", result.Error)
		}

		// Store the generated UUID for the artist
		artistUUIDs[artist.Name] = artist.UUID

		if result.Error != nil {
			log.Fatalf("Failed to seed stage data: %v", result.Error)
		}
		var randomStage models.Stage
		db.Order("RANDOM()").First(&randomStage)

		weekNumberValue := gofakeit.Number(1, 52)
		weekNumber := &weekNumberValue

		// Generate fake data
		yearValue := gofakeit.Year()
		dayValue := gofakeit.WeekDay()
		startTimeValue := gofakeit.Date()
		endTimeValue := gofakeit.Date()

		// Create pointers to the values
		year := &yearValue
		day := &dayValue
		startTime := &startTimeValue
		endTime := &endTimeValue
		// Generate fake data for TimetableEntry
		timetableEntry := models.TimetableEntry{
			ID:         uuid.New(),
			StageID:    randomStage.ID,
			ArtistID:   artist.UUID,
			WeekNumber: weekNumber,
			Year:       year,
			Day:        day,
			StartTime:  startTime,
			EndTime:    endTime,
		}
		result = db.Create(&timetableEntry)
		if result.Error != nil {
			log.Fatalf("Failed to seed timetable entry data: %v", result.Error)
		}
	}

	// Now you have artistUUIDs and stageUUIDs maps that store the generated UUIDs
	// for Artists and Stages, which you can use for other seeding tasks if needed.
}
