package artistApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"gorm.io/gorm"
)

const (
	tokenURL     = "https://api.soundcloud.com/oauth2/token" // Token endpoint
	clientID     = "tSJvk12uBEeL3YjyxxQ4IDloTcloGJzl"
	clientSecret = "JU9HgHJ7VXJgu4lzFH6BqscnqUOwQj7D"
)

type OAuthToken struct {
	gorm.Model
	AccessToken        string    `gorm:"type:varchar(255);not null"`
	RefreshToken       string    `gorm:"type:varchar(255);not null"`
	AccessTokenExpiry  time.Time `gorm:"not null"`
	RefreshTokenExpiry time.Time
}

// AccessTokenResponse represents the response from the OAuth server
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// TableName overrides the table name used by GORM
func (OAuthToken) TableName() string {
	return "oauth_tokens"
}

// RequestAccessToken requests an access token using client credentials
func RequestAccessToken() (*AccessTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func RefreshAccessToken(clientID, clientSecret, refreshToken string) (*AccessTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func StartTokenRefreshScheduler(db *gorm.DB) {
	ticker := time.NewTicker(30 * time.Minute) // Adjust the duration according to your needs
	go func() {
		for {
			select {
			case <-ticker.C:
				refreshTokenIfNeeded(db, clientID, clientSecret)
			}
		}
	}()
}

func refreshTokenIfNeeded(db *gorm.DB, clientID, clientSecret string) {
	// Retrieve the latest token from the database
	token, err := GetLatestToken(db)
	if err != nil {
		log.Printf("Error retrieving token: %v", err)
		return
	}

	// Check if the token is close to expiry
	if time.Now().Add(10 * time.Minute).After(token.AccessTokenExpiry) { // 10 minutes before expiry
		log.Printf("Asking for new token")
		refreshedToken, err := RefreshAccessToken(clientID, clientSecret, token.RefreshToken)
		if err != nil {
			log.Printf("Error refreshing token: %v", err)
			return
		}

		log.Printf("Save new token for new token: %v", refreshedToken.AccessToken)
		// Save the refreshed token
		result, err := SaveToken(db, refreshedToken)
		if err != nil {
			log.Printf("Error saving refreshed token: %v", err)
		}
		log.Printf("Saved new token: %v", result.AccessToken)
	}
}

func SaveToken(db *gorm.DB, token *AccessTokenResponse) (*OAuthToken, error) {
	expiryTime := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	// Assuming you have a way to identify the existing token record (e.g., a single record or based on client ID)
	var existingToken OAuthToken
	result := db.First(&existingToken)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Insert new token if no existing token is found
		newToken := OAuthToken{
			AccessToken:        token.AccessToken,
			RefreshToken:       token.RefreshToken,
			AccessTokenExpiry:  expiryTime,
			RefreshTokenExpiry: expiryTime,
		}
		result = db.Create(&newToken)
		if result.Error != nil {
			return nil, result.Error
		}
		return &newToken, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	// Update existing token
	existingToken.AccessToken = token.AccessToken
	existingToken.RefreshToken = token.RefreshToken
	existingToken.AccessTokenExpiry = expiryTime
	existingToken.RefreshTokenExpiry = expiryTime

	result = db.Save(&existingToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &existingToken, nil
}

func GetLatestToken(db *gorm.DB) (*OAuthToken, error) {
	var token OAuthToken
	result := db.Order("access_token_expiry DESC").First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// No token record found
			return nil, nil
		}
		return nil, result.Error
	}
	return &token, nil
}
