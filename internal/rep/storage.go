package rep

import (
	"context"

	"gorm.io/gorm"
)

type Storage interface {
	GetGormDB(ctx context.Context) *gorm.DB
}
