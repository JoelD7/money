package savings

import (
	"context"
	"github.com/JoelD7/money/backend/models"
)

type Repository interface {
	GetSavings(ctx context.Context, email string) ([]*models.Saving, error)
	CreateSaving(ctx context.Context, saving *models.Saving) error
}
