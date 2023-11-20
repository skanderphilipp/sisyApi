package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
	"github.com/skanderphilipp/sisyApi/internal/domain/event"
	"github.com/skanderphilipp/sisyApi/internal/domain/stage"
	"github.com/skanderphilipp/sisyApi/internal/domain/venue"
	"github.com/skanderphilipp/sisyApi/internal/infrastructure/api/artistApi"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProvideDatabase() (*gorm.DB, error) {
	// Read environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Create connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	err = db.AutoMigrate(&artist.Artist{}, &artist.SocialMedia{}, &venue.Venue{}, &stage.Stage{}, &event.Event{}, &event.TimetableEntry{}, &artistApi.OAuthToken{})

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db, nil
}

// Ensure you call this function from somewhere in your application, typically main.go
func CloseDatabaseConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error on getting DB from GORM: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Error on closing db connection: %v", err)
	}
}
