package repository

import "gorm.io/gorm"

type StageRepository struct {
	db *gorm.DB
}

func NewStageRepository(db *gorm.DB) *StageRepository {
	return &StageRepository{db: db}
}

func getByVenueID() {

}
