package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/blnto/blnto_service/internal"
	"github.com/blnto/blnto_service/internal/domain/artist"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type ArtistData struct {
	Name           string `json:"name"`
	SoundcloudLink string `json:"soundcloudLink,omitempty"`
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Define flags
	importFilePath := flag.String("file", "", "/data/artistdata.json")
	flag.Parse()

	// Validate the input
	if *importFilePath == "" {
		log.Fatal("You must specify a file path using the -file flag.")
	}

	// Initialize database connection
	app, err := internal.InitializeDependencies()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Load artist data from the specified file
	artists, err := loadArtistsFromFile(*importFilePath)
	if err != nil {
		log.Fatalf("Failed to load artists from file: %v", err)
	}

	// Perform the mass import
	err = massImportArtists(app.DB, artists)
	if err != nil {
		log.Fatalf("Failed to import artists: %v", err)
	}

	fmt.Println("Artists imported successfully.")
}

func loadArtistsFromFile(filePath string) ([]artist.Artist, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var artistData []ArtistData
	err = json.Unmarshal(fileData, &artistData)
	if err != nil {
		return nil, err
	}

	var artists []artist.Artist
	for _, ad := range artistData {
		var artistEntry artist.Artist
		artistEntry.Name = ad.Name

		if ad.SoundcloudLink != "" {
			artistEntry.SCPermalink = new(string)
			*artistEntry.SCPermalink = ad.SoundcloudLink
		}

		artists = append(artists, artistEntry)
	}

	return artists, nil
}

func massImportArtists(db *gorm.DB, artists []artist.Artist) error {
	for _, artistData := range artists {
		fmt.Println("Artist import: ", artistData.Name)
		var existingArtist artist.Artist
		result := db.Where("name = ?", artistData.Name).First(&existingArtist)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("Artist does not exist, create a new record")
			// Artist does not exist, create a new record
			if err := db.Create(&artistData).Error; err != nil {
				return err
			}
		} else {
			// Artist already exists, do nothing
			fmt.Println("Artist already exists, skipping")
		}
	}
	return nil
}
