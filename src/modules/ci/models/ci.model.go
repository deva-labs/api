package ci

import (
	projects "deva/src/modules/projects/models"
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CiPipeline struct {
	ID            uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID     uuid.UUID        `gorm:"type:uuid;not null"`
	Project       projects.Project `gorm:"foreignKey:ProjectID;references:ID"`
	Provider      string           `gorm:"not null"`
	Status        string           `gorm:"not null;default:'pending'"`
	CommitSHA     string
	LogURL        string
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedByUser users.User     `gorm:"foreignKey:UpdatedBy;references:ID"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateCIPipelines(db *gorm.DB) error {
	return db.AutoMigrate(&CiPipeline{})
}
