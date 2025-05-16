package teams

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	users "skypipe/src/modules/users/models"
	"time"
)

type TeamMember struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TeamID        uuid.UUID      `gorm:"type:uuid;not null"`
	Team          Team           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:TeamID;references:ID;"`
	UserID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	User          users.User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:TeamID;references:ID;"`
	Role          string         `gorm:"not null"`
	JoinedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedByUser users.User     `gorm:"foreignKey:UpdatedBy;references:ID;"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateTeamMembers(db *gorm.DB) error {
	return db.AutoMigrate(&TeamMember{})
}
