package templates

import (
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ProjectTemplate struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name          string    `gorm:"not null;unique"`
	Description   string
	Language      string         `gorm:"not null"`
	IsOfficial    bool           `gorm:"not null;default:false"`
	UpdatedBy     uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedByUser users.User     `gorm:"foreignKey:UpdatedBy;references:ID"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func MigrateProjectTemplates(db *gorm.DB) error {
	return db.AutoMigrate(&ProjectTemplate{})
}
