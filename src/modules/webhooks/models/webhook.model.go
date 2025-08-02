package webhooks

import (
	projects "deva/src/modules/projects/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Webhook struct {
	ID        uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID uuid.UUID        `gorm:"type:uuid;not null"`
	Project   projects.Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ProjectID;references:ID;"`
	URL       string           `gorm:"not null"`
	EventType string           `gorm:"not null"` // "deployment.success", "build.failure", etc.
	Secret    string           `gorm:"not null"` // HMAC secret
	IsActive  bool             `gorm:"not null;default:true"`
	CreatedAt time.Time        `gorm:"autoCreateTime"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt   `gorm:"index"`
}

func MigrateWebhooks(db *gorm.DB) error {
	return db.AutoMigrate(&Webhook{})
}
