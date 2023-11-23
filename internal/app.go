package internal

import (
	"github.com/blnto/blnto_service/internal/api/graphql/resolvers"
	"github.com/blnto/blnto_service/internal/application/service"
	"github.com/blnto/blnto_service/internal/infrastructure/repository"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"os"
)

type App struct {
	DB               *gorm.DB
	ArtistService    *service.ArtistService
	EventService     *service.EventService
	StageService     *service.StageService
	VenueService     *service.VenueService
	Logger           *zap.Logger
	Loggerfile       *os.File
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
		Loggerfile:       config.Loggerfile,
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
	logger, file, err := provideLogger()

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
		Loggerfile:       file,
		Resolver:         resolver,
		ArtistRepository: artistRepo,
		VenueRepository:  venueRepo,
		EventRepository:  eventRepo,
		StageRepository:  stageRepo,
	}
	return NewApp(appConfig), nil
}

func provideLogger() (*zap.Logger, *os.File, error) {
	// Create a file to write logs to
	file, err := os.OpenFile("logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create a core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zapcore.InfoLevel,
	)

	// Create a logger
	logger := zap.New(core)

	// Example log
	logger.Info("This is a JSON log message", zap.String("type", "example"))
	return logger, file, nil
}
