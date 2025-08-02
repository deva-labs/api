package key_token

import (
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AccessToken struct {
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

// BeforeCreate sets the default expiration time for the access token.
func (a *AccessToken) BeforeCreate(tx *gorm.DB) (err error) {
	a.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // 7 days from now
	return
}

// MigrateAccessTokens migrates the AccessToken models.
func MigrateAccessTokens(db *gorm.DB) error {
	return db.AutoMigrate(&AccessToken{})
}
