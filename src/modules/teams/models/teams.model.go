package teams

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Team struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string         `gorm:"not null"`
	OwnerID   uuid.UUID      `gorm:"type:uuid;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateTeams(db *gorm.DB) error {
	return db.AutoMigrate(&Team{})
}
