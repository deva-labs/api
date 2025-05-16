package users

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	roles "skypipe/src/modules/roles/models"
	"time"
)

type UserRole struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null"`
	RoleID        uuid.UUID      `gorm:"type:uuid;not null"`
	User          User           `gorm:"foreignKey:UserID;references:ID"`
	Role          roles.Role     `gorm:"foreignKey:RoleID;references:ID"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedByUser User           `gorm:"foreignKey:UpdatedBy;references:ID"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateUserRole(db *gorm.DB) error {
	return db.AutoMigrate(&UserRole{})
}
