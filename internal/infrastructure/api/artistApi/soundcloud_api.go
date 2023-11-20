package artistApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	baseAPIURL = "https://api.soundcloud.com"
)

// SoundCloudClient holds the configuration for the API client
type SoundCloudClient struct {
	DB           *gorm.DB
	ClientID     string
	clientSecret string
}

// Artist represents the artist data structure for SoundCloud
type SCArtist struct {
	ID          int    `json:"id"`
	City        string `json:"city"`
	AvatarURL   string `json:"avatar_url"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Country     string `json:"country"`
	Permalink   string `json:"permalink"`
}

type SCTrack struct {
	HTTPMp3URL    string `json:"http_mp3_128_url"`
	HLSMp3URL     string `json:"hls_mp3_128_url"`
	HLSOpusURL    string `json:"hls_opus_64_url"`
	PreviewMp3URL string `json:"preview_mp3_128_url"`
}

func NewClient(db *gorm.DB) *SoundCloudClient {
	return &SoundCloudClient{DB: db}
}

// FetchArtistByPermalink fetches an artist by their permalink
func (sc *SoundCloudClient) FetchArtistByPermalink(permalink string) (*SCArtist, error) {
	token, err := GetLatestToken(sc.DB)

	if err != nil {
		return nil, err
	}
	// Check if the token is close to expiry and refresh if necessary
	if time.Now().After(token.AccessTokenExpiry) {

		refreshedToken, refreshErr := RefreshAccessToken(clientID, clientSecret, token.RefreshToken)
		if refreshErr != nil {
			return nil, err
		}

		token, err = SaveToken(sc.DB, refreshedToken)

		if err != nil {
			zap.L().Error("failed to save token", zap.Error(err))
			return nil, err
		}
	}

	// Construct the request URL
	resolveURL := fmt.Sprintf("https://api.soundcloud.com/resolve?url=%s", url.QueryEscape(permalink))
	// Create an HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", resolveURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token.AccessToken))

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Decode the response
	var artist SCArtist
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &artist, nil
}

// FetchArtistById fetches an artist by their SoundCloud ID
func (sc *SoundCloudClient) FetchArtistById(id string) (*SCArtist, error) {
	// Implement the API call here
	token, err := GetLatestToken(sc.DB)

	if err != nil {
		return nil, err
	}
	// Check if the token is close to expiry and refresh if necessary
	if time.Now().After(token.AccessTokenExpiry) {
		refreshedToken, refreshErr := RefreshAccessToken(clientID, clientSecret, token.RefreshToken)
		if refreshErr != nil {
			return nil, err
		}

		token, err = SaveToken(sc.DB, refreshedToken)

		if err != nil {
			zap.L().Error("failed to save token", zap.Error(err))
			return nil, err
		}
	}

	// Construct the request URL
	resolveURL := fmt.Sprintf("https://api.soundcloud.com/users/%s", id)
	// Create an HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", resolveURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", token.AccessToken))

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Decode the response
	var artist SCArtist
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &artist, nil
}

// FetchTracksByArtistId fetches tracks by the artist's SoundCloud ID
func (sc *SoundCloudClient) FetchTracksByArtistId(artistId string) ([]SCTrack, error) {
	return nil, nil
	// Implement the API call here
}
