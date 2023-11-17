package internal

import (
	"github.com/skanderphilipp/sisyApi/internal/api/graphql/resolvers"
	"github.com/skanderphilipp/sisyApi/internal/application/service"
	"github.com/skanderphilipp/sisyApi/internal/infrastructure/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	DB               *gorm.DB
	ArtistService    *service.ArtistService
	Logger           *zap.Logger
	Resolver         *resolvers.Resolver
	ArtistRepository *repository.ArtistRepository
	VenueRepository  *repository.VenueRepository
	EventRepository  *repository.EventRepository
	// other components like service clients, configuration, etc.
}

func NewApp(db *gorm.DB, artistService *service.ArtistService, logger *zap.Logger, resolver *resolvers.Resolver, artistRepository *repository.ArtistRepository, venueRepository *repository.VenueRepository, eventRepository *repository.EventRepository) *App {
	return &App{
		DB:               db,
		ArtistService:    artistService,
		Logger:           logger,
		Resolver:         resolver,
		ArtistRepository: artistRepository,
		VenueRepository:  venueRepository,
		EventRepository:  eventRepository,
	}
}

func InitializeDependencies() (*App, error) {
	// Create a database connection
	db, err := repository.ProvideDatabase()
	if err != nil {
		return nil, err
	}

	// Create a repository
	artistRepo := repository.NewArtistRepository(db)
	venueRepo := repository.NewVenueRepository(db)
	eventRepo := repository.NewEventRepository(db)
	// Create a service
	artistService := service.NewArtistService(artistRepo)

	// Create a logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	// Create a resolver
	resolver := resolvers.NewResolver(artistService, artistRepo, venueRepo, eventRepo)

	return NewApp(db, artistService, logger, resolver, artistRepo, venueRepo, eventRepo), nil
}

func provideLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	return logger, err
}
func provideAritstService(repo *repository.ArtistRepository) *service.ArtistService {
	return service.NewArtistService(repo)
}

func provideArtistRepository(db *gorm.DB) *repository.ArtistRepository {
	return repository.NewArtistRepository(db)
}
