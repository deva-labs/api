package projects

import (
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ProjectFile struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID     uuid.UUID      `gorm:"type:uuid;not null"`
	Project       Project        `gorm:"foreignKey:ProjectID;references:ID"`
	Path          string         `gorm:"not null"`
	Content       string         `gorm:"type:text"`
	IsGenerated   bool           `gorm:"not null;default:false"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedByUser users.User     `gorm:"foreignKey:UpdatedBy;references:ID"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateProjectFiles(db *gorm.DB) error {
	return db.AutoMigrate(&ProjectFile{})
}
