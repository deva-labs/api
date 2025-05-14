package projects

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ProjectFile struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;not null"`
	Path        string         `gorm:"not null"`
	Content     string         `gorm:"type:text"`
	IsGenerated bool           `gorm:"not null;default:false"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func MigrateProjectFiles(db *gorm.DB) error {
	return db.AutoMigrate(&ProjectFile{})
}
