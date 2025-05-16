package deployments

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	users "skypipe/src/modules/users/models"
	"time"
)

type DeploymentTarget struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TeamID        uuid.UUID      `gorm:"type:uuid;default:null"` // Need change to direct team id after MVP
	Name          string         `gorm:"not null"`
	Type          string         `gorm:"not null"` // "kubernetes", "docker", etc.
	Host          string         `gorm:"not null"`
	Auth          string         `gorm:"type:jsonb;not null"` // Encrypted auth credentials
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;default:null"`
	UpdatedByUser users.User     `gorm:"foreignKey:UpdatedBy;references:ID"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateDeploymentTargets(db *gorm.DB) error {
	return db.AutoMigrate(&DeploymentTarget{})
}
