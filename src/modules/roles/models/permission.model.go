package roles

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Resource  string         `gorm:"not null"` // ex: "project", "deployment"
	Action    string         `gorm:"not null"` // ex: "create", "read", "update", "delete"
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigratePermission(db *gorm.DB) error {
	return db.AutoMigrate(&Permission{})
}
