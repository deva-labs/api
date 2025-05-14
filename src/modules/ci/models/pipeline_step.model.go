package ci

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PipelineStep struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PipelineID  uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"not null"`
	Status      string    `gorm:"not null;default:'pending'"`
	Log         string    `gorm:"type:text"`
	StartedAt   time.Time
	CompletedAt time.Time
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func MigratePipelineSteps(db *gorm.DB) error {
	return db.AutoMigrate(&PipelineStep{})
}
