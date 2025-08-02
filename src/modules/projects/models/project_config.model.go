package projects

import (
	deployments "deva/src/modules/deployments/models"
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ProjectConfig struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Project   Project   `gorm:"foreignKey:ProjectID;references:ID"`

	UpdatedBy     uuid.UUID  `gorm:"type:uuid;not null"`
	UpdatedByUser users.User `gorm:"foreignKey:UpdatedBy;references:ID"`

	Language       string `gorm:"not null"`
	Framework      string
	EnvVars        string `gorm:"type:jsonb"`
	CITool         string
	DeployTargetID uuid.UUID                    `gorm:"type:uuid;not null"`
	DeployTarget   deployments.DeploymentTarget `gorm:"foreignKey:DeployTargetID;references:ID"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateProjectConfigs(db *gorm.DB) error {
	return db.AutoMigrate(&ProjectConfig{})
}
