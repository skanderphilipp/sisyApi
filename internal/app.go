package internal

import (
	"github.com/skanderphilipp/sisyApi/internal/application/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	DB            *gorm.DB
	ArtistService *service.ArtistService
	Logger        *zap.Logger
	// other components like service clients, configuration, etc.
}

func NewApp(db *gorm.DB, artistService *service.ArtistService, logger *zap.Logger) *App {
	return &App{
		DB:            db,
		ArtistService: artistService,
		Logger:        logger,
		// Initialize other components
	}
}
