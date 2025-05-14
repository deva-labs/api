package notifications

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Type      string    `gorm:"not null"` // "deployment", "billing", "invite", etc.
	Message   string    `gorm:"not null"`
	IsRead    bool      `gorm:"not null;default:false"`
	Metadata  string    `gorm:"type:jsonb"` // Additional context
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ReadAt    time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateNotifications(db *gorm.DB) error {
	return db.AutoMigrate(&Notification{})
}
