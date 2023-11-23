package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blnto/blnto_service/internal"
	"github.com/blnto/blnto_service/internal/infrastructure/api/artistApi"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize database connection
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	app, err := internal.InitializeDependencies() // Adjust with your actual setup function
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Fetch artists with permalink
	artists, err := app.ArtistRepository.FindAllWithPermalink(context.TODO())
	if err != nil {
		log.Fatalf("Failed to fetch artists: %v", err)
	}

	// Initialize SoundCloud API client
	scClient := artistApi.NewClient(app.DB) // Adjust accordingly

	for _, artist := range artists {
		// Fetch artist info from SoundCloud
		permalink := artist.SCPermalink
		scInfo, err := scClient.FetchArtistByPermalink(*permalink)
		if err != nil {
			log.Printf("Failed to fetch SoundCloud info for %s: %v", artist.Name, err)
			continue
		}

		// Update artist with SoundCloud info
		artist.UpdateWithSoundCloudInfo(*scInfo) // Implement this method on your Artist model

		// Save updated artist to the database
		if artist, err := app.ArtistRepository.Update(context.TODO(), &artist); err != nil {
			log.Printf("Failed to update artist %s: %v", artist.Name, err)
		}
	}

	fmt.Println("Artist information updated successfully.")
}
