package captcha

import (
	users "deva/src/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Captcha struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Text      string         `gorm:"type:text;not null"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null"`
	User      users.User     `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateCaptcha(db *gorm.DB) error {
	return db.AutoMigrate(&Captcha{})
}
