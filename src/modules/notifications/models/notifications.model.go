package notifications

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	users "skypipe/src/modules/users/models"
	"time"
)

type Notification struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null"`
	User          users.User `gorm:"foreignKey:UserID;references:ID"`
	Type          string     `gorm:"not null"` // "deployment", "billing", "invite", etc.
	Message       string     `gorm:"not null"`
	IsRead        bool       `gorm:"not null;default:false"`
	Metadata      string     `gorm:"type:jsonb"` // Additional context
	UpdatedBy     uuid.UUID  `gorm:"type:uuid;not null"`
	UpdatedByUser users.User `gorm:"foreignKey:UpdatedBy;references:ID"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	ReadAt        time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateNotifications(db *gorm.DB) error {
	return db.AutoMigrate(&Notification{})
}
