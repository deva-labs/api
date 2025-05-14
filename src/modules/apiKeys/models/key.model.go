package apiKeys

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type APIKey struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"not null"`
	TokenHash string    `gorm:"not null;uniqueIndex"`
	Scopes    string    `gorm:"type:jsonb;not null"` // Array of permissions
	LastUsed  time.Time
	ExpiresAt time.Time
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateAPIKeys(db *gorm.DB) error {
	return db.AutoMigrate(&APIKey{})
}
