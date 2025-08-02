package logs

import (
	projects "deva/src/modules/projects/models"
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ActivityLog struct {
	ID        uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID        `gorm:"type:uuid;not null"`
	User      users.User       `gorm:"foreignKey:UserID;references:ID"`
	ProjectID uuid.UUID        `gorm:"type:uuid"`
	Project   projects.Project `gorm:"foreignKey:ProjectID;references:ID"`
	Action    string           `gorm:"not null"` // e.g., "project.create", "deployment.start"
	Metadata  string           `gorm:"type:jsonb"`
	IPAddress string           `gorm:"not null"`
	CreatedAt time.Time        `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt   `gorm:"index"`
}

func MigrateActivityLogs(db *gorm.DB) error {
	return db.AutoMigrate(&ActivityLog{})
}
