//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/skanderphilipp/sisyApi/internal"
	"github.com/skanderphilipp/sisyApi/internal/application/service"
	"github.com/skanderphilipp/sisyApi/internal/domain/artist"
	"github.com/skanderphilipp/sisyApi/internal/infrastructure/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitializeApp() (*internal.App, error) {
	wire.Build(
		repository.ProvideDatabase,
		provideArtistRepository,
		provideAritstService,
		provideLogger,
		internal.NewApp,
	)

	return &internal.App{}, nil
}

func provideLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	return logger, err
}
func provideAritstService(repo artist.Repository) *service.ArtistService {
	return service.NewArtistService(repo)
}

func provideArtistRepository(db *gorm.DB) artist.Repository {
	return repository.NewArtistRepository(db)
}
