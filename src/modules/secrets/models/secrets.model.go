package secrets

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Secret struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID      uuid.UUID      `gorm:"type:uuid;not null"`
	Key            string         `gorm:"not null"`
	ValueEncrypted string         `gorm:"not null"` // Encrypted value
	Scope          string         `gorm:"not null"` // "build", "runtime", "both"
	Version        int            `gorm:"not null;default:1"`
	CreatedBy      uuid.UUID      `gorm:"type:uuid;not null"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func MigrateSecrets(db *gorm.DB) error {
	return db.AutoMigrate(&Secret{})
}
