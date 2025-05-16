package apiKeys

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	users "skypipe/src/modules/users/models"
	"time"
)

type APIKey struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	User          users.User `gorm:"foreignKey:UserID;references:ID"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null"`
	UpdatedBy     uuid.UUID  `gorm:"type:uuid;not null"`
	Name          string     `gorm:"not null"`
	TokenHash     string     `gorm:"not null;uniqueIndex"`
	Scopes        []string   `gorm:"type:jsonb;not null"`
	LastUsed      time.Time
	ExpiresAt     time.Time
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	UpdatedByUser users.User     `gorm:"foreignKey:UpdatedBy;references:ID"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateAPIKeys(db *gorm.DB) error {
	return db.AutoMigrate(&APIKey{})
}
