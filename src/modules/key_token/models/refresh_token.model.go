package key_token

import (
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID      `gorm:"not null"`
	User          users.User     `gorm:"foreignKey:UserID;references:ID"`
	Token         string         `gorm:"type:varchar(256);not null"`
	Status        bool           `gorm:"default:true"` // true = active, false = revoked
	ExpiresAt     time.Time      `gorm:"not null"`     // Token expiration time
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedByUser users.User     `gorm:"foreignKey:UpdatedBy;references:ID"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate sets the default expiration time for the refresh token.
func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	r.ExpiresAt = time.Now().Add(30 * 24 * time.Hour) // 30 days from now
	return
}

// MigrateRefreshTokens migrates the RefreshToken models.
func MigrateRefreshTokens(db *gorm.DB) error {
	return db.AutoMigrate(&RefreshToken{})
}
