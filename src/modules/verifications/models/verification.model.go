package verifications

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type VerificationCode struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email      string    `gorm:"type:varchar(255);not null"`
	Code       string    `gorm:"type:varchar(255);not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	InputCount int       `gorm:"type:int;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime;"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;"`
	DeletedAt  time.Time `gorm:"index;"`
}

func MigrateVerificationCode(db *gorm.DB) error {
	return db.AutoMigrate(&VerificationCode{})
}
