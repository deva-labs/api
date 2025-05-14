package logs

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ActivityLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ProjectID uuid.UUID `gorm:"type:uuid"`
	Action    string    `gorm:"not null"` // e.g., "project.create", "deployment.start"
	Metadata  string    `gorm:"type:jsonb"`
	IPAddress string
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateActivityLogs(db *gorm.DB) error {
	return db.AutoMigrate(&ActivityLog{})
}
