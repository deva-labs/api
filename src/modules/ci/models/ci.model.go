package ci

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CIPipeline struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null"`
	Provider  string    `gorm:"not null"`
	Status    string    `gorm:"not null;default:'pending'"`
	CommitSHA string
	LogURL    string
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateCIPipelines(db *gorm.DB) error {
	return db.AutoMigrate(&CIPipeline{})
}
