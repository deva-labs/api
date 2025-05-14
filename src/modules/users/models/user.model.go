package users

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Name      string         `gorm:"not null"`
	Role      string         `gorm:"not null;default:'user'"`
	PlanID    *uuid.UUID     // Nullable foreign key
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateUserCore(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
