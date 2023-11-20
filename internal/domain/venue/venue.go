package venue

import (
	"time"

	"github.com/google/uuid"
	"github.com/skanderphilipp/sisyApi/internal/domain/stage"
	"gorm.io/gorm"
)

type Venue struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description *string        `json:"description,omitempty"`
	Stages      []*stage.Stage `gorm:"foreignKey:VenueID" json:"stages,omitempty"`
	//gorm additinonal fields
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// Additional fields like CreatedAt, UpdatedAt can be added.
}

// Venue BeforeCreate hook
func (v *Venue) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return
}
