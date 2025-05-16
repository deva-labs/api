package roles

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type RolePermission struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoleID       uuid.UUID      `gorm:"type:uuid;not null"`
	PermissionID uuid.UUID      `gorm:"type:uuid;not null"`
	Role         Role           `gorm:"foreignKey:RoleID;references:ID"`
	Permission   Permission     `gorm:"foreignKey:PermissionID;references:ID"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func MigrateRolePermission(db *gorm.DB) error {
	return db.AutoMigrate(&RolePermission{})
}
