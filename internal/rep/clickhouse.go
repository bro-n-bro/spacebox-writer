package rep

import (
	"context"

	"gorm.io/gorm"
)

type Clickhouse interface {
	GetGormDB(ctx context.Context) *gorm.DB
}
