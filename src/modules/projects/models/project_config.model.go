package projects

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ProjectConfig struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Language     string    `gorm:"not null"`
	Framework    string
	Port         int
	EnvVars      string `gorm:"type:jsonb"`
	CITool       string
	DeployTarget string
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func MigrateProjectConfigs(db *gorm.DB) error {
	return db.AutoMigrate(&ProjectConfig{})
}
