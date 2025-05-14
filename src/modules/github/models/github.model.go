package github

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type GitHubIntegration struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	GithubUserID string    `gorm:"not null"`
	AccessToken  string    `gorm:"not null"` // Encrypted token
	RefreshToken string    // Encrypted token
	ExpiresAt    time.Time
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func MigrateGitHubIntegrations(db *gorm.DB) error {
	return db.AutoMigrate(&GitHubIntegration{})
}
