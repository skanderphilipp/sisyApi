package artist

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Artist, error)
	// Other necessary methods like Save, Delete, FindAll...
}
