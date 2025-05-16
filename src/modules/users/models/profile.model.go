package users

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	User          User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	FullName      string    `gorm:"type:varchar(100)"`
	Bio           string    `gorm:"type:text"`
	AvatarURL     string    `gorm:"type:text"`
	Gender        string    `gorm:"type:varchar(20)"`
	Birthday      *time.Time
	Phone         string         `gorm:"type:varchar(20);uniqueIndex;not null"`
	Address       string         `gorm:"type:text"`
	Country       string         `gorm:"type:varchar(100)"`
	City          string         `gorm:"type:varchar(100)"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedByUser User           `gorm:"foreignKey:UpdatedBy;references:ID;onUpdate:CASCADE;onDelete:CASCADE;"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateProfile(db *gorm.DB) error {
	return db.AutoMigrate(&Profile{})
}
