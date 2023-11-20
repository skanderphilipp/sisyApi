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
	EventService     *service.EventService
	StageService     *service.StageService
	VenueService     *service.VenueService
	Logger           *zap.Logger
	Resolver         *resolvers.Resolver
	ArtistRepository *repository.ArtistRepository
	VenueRepository  *repository.VenueRepository
	EventRepository  *repository.EventRepository
	StageRepository  *repository.StageRepository
}

func NewApp(config *App) *App {
	return &App{
		DB:               config.DB,
		ArtistService:    config.ArtistService,
		EventService:     config.EventService,
		StageService:     config.StageService,
		VenueService:     config.VenueService,
		Logger:           config.Logger,
		Resolver:         config.Resolver,
		ArtistRepository: config.ArtistRepository,
		VenueRepository:  config.VenueRepository,
		EventRepository:  config.EventRepository,
		StageRepository:  config.StageRepository,
	}
}

func InitializeDependencies() (*App, error) {
	// Create a database connection
	db, err := repository.ProvideDatabase()

	// seed.SeedDatabase(db)

	if err != nil {
		return nil, err
	}

	// Create a repository
	artistRepo := repository.NewArtistRepository(db)
	venueRepo := repository.NewVenueRepository(db)
	eventRepo := repository.NewEventRepository(db)
	stageRepo := repository.NewStageRepository(db)
	// Create a service
	artistService := service.NewArtistService(artistRepo)
	eventService := service.NewEventService(eventRepo)
	stageService := service.NewStageService(stageRepo)
	venueService := service.NewVenueService(venueRepo)

	// Create a logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	// Create a resolver
	resolver := resolvers.NewResolver(artistService, eventService, stageService, venueService)

	appConfig := &App{
		DB:               db,
		ArtistService:    artistService,
		EventService:     eventService,
		StageService:     stageService,
		VenueService:     venueService,
		Logger:           logger,
		Resolver:         resolver,
		ArtistRepository: artistRepo,
		VenueRepository:  venueRepo,
		EventRepository:  eventRepo,
		StageRepository:  stageRepo,
	}
	return NewApp(appConfig), nil
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
