package templates

import (
	projects "deva/src/modules/projects/models"
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UsageMetric struct {
	ID            uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID        `gorm:"type:uuid;not null"`
	User          users.User       `gorm:"foreignKey:UserID;references:ID"`
	ProjectID     uuid.UUID        `gorm:"type:uuid"`
	Project       projects.Project `gorm:"foreignKey:ProjectID;references:ID"`
	MetricType    string           `gorm:"not null"` // "deployment", "build", "storage", etc.
	Count         int              `gorm:"not null"`
	PeriodStart   time.Time        `gorm:"not null"`
	PeriodEnd     time.Time        `gorm:"not null"`
	UpdatedBy     uuid.UUID        `gorm:"type:uuid;not null"`
	UpdatedByUser users.User       `gorm:"foreignKey:UpdatedBy;references:ID"`
	CreatedAt     time.Time        `gorm:"autoCreateTime"`
	UpdatedAt     time.Time        `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt   `gorm:"index"`
}

func MigrateUsageMetrics(db *gorm.DB) error {
	return db.AutoMigrate(&UsageMetric{})
}
