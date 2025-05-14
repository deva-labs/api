package projects

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TeamID     uuid.UUID `gorm:"type:uuid;not null"`
	Name       string    `gorm:"not null"`
	RepoURL    string
	SourceType string         `gorm:"not null"`
	Status     string         `gorm:"not null;default:'active'"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func MigrateProjects(db *gorm.DB) error {
	return db.AutoMigrate(&Project{})
}
