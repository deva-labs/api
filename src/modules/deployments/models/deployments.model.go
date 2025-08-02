package deployments

import (
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Deployment struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Platform    string     `gorm:"not null"`
	Status      string     `gorm:"not null;default:'pending'"`
	Log         string     `gorm:"type:text"`
	TriggeredBy uuid.UUID  `gorm:"type:uuid;not null"`
	User        users.User `gorm:"foreignKey:TriggeredBy;references:ID"`
	StartedAt   time.Time
	FinishedAt  time.Time
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func MigrateDeployments(db *gorm.DB) error {
	return db.AutoMigrate(&Deployment{})
}
