package github

import (
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type GithubIntegration struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex"`
	User          users.User `gorm:"foreignKey:UserID;references:ID"`
	GithubID      string     `gorm:"not null"`
	AccessToken   string     `gorm:"not null"` // Encrypted token
	RefreshToken  string     `gorm:"not null"` // Encrypted token
	UpdatedBy     uuid.UUID  `gorm:"type:uuid;not null"`
	UpdatedByUser users.User `gorm:"foreignKey:UpdatedBy;references:ID"`
	ExpiresAt     time.Time
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateGitHubIntegrations(db *gorm.DB) error {
	return db.AutoMigrate(&GithubIntegration{})
}
