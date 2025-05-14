package deployments

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Deployment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID   uuid.UUID `gorm:"type:uuid;not null"`
	Platform    string    `gorm:"not null"`
	Status      string    `gorm:"not null;default:'pending'"`
	Log         string    `gorm:"type:text"`
	TriggeredBy uuid.UUID `gorm:"type:uuid;not null"`
	StartedAt   time.Time
	FinishedAt  time.Time
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func MigrateDeployments(db *gorm.DB) error {
	return db.AutoMigrate(&Deployment{})
}
