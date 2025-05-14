package teams

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TeamMember struct {
	TeamID    uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Role      string         `gorm:"not null"`
	JoinedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateTeamMembers(db *gorm.DB) error {
	return db.AutoMigrate(&TeamMember{})
}
