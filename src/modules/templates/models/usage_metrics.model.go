package templates

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UsageMetric struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	ProjectID   uuid.UUID `gorm:"type:uuid"`
	MetricType  string    `gorm:"not null"` // "deployment", "build", "storage", etc.
	Count       int       `gorm:"not null"`
	PeriodStart time.Time `gorm:"not null"`
	PeriodEnd   time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func MigrateUsageMetrics(db *gorm.DB) error {
	return db.AutoMigrate(&UsageMetric{})
}
