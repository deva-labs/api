package deployments

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type DeploymentTarget struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TeamID    uuid.UUID      `gorm:"type:uuid;not null"`
	Name      string         `gorm:"not null"`
	Type      string         `gorm:"not null"` // "kubernetes", "docker", etc.
	Host      string         `gorm:"not null"`
	Auth      string         `gorm:"type:jsonb;not null"` // Encrypted auth credentials
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateDeploymentTargets(db *gorm.DB) error {
	return db.AutoMigrate(&DeploymentTarget{})
}
