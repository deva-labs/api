package plans

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Plan struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name                 string         `gorm:"not null;unique"`
	Price                float64        `gorm:"not null"`
	MaxProjects          int            `gorm:"not null"`
	MaxDeploymentsPerDay int            `gorm:"not null"`
	CreatedAt            time.Time      `gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

func MigratePlan(db *gorm.DB) error {
	return db.AutoMigrate(&Plan{})
}
