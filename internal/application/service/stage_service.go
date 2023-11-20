package service

import (
	"github.com/skanderphilipp/sisyApi/internal/domain/models"
	"github.com/skanderphilipp/sisyApi/internal/domain/stage"
	"github.com/skanderphilipp/sisyApi/internal/infrastructure/repository"
)

type StageService struct {
	repo *repository.StageRepository
}

func NewStageService(repo *repository.StageRepository) *StageService {
	return &StageService{repo: repo}
}

func mapGormStageToGqlStage(gormStage *stage.Stage) *models.Stage {
	gqlStage := &models.Stage{
		ID:   gormStage.ID,
		Name: gormStage.StageName,
	}

	return gqlStage
}

func mapGqlStageToGormStage(gqlStage *models.Stage) *stage.Stage {
	gormStage := &stage.Stage{
		ID:        gqlStage.ID,
		StageName: gqlStage.Name,
	}

	return gormStage
}
