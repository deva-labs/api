package users

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	plans "skypipe/src/modules/plans/models"
	"time"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Name      string         `gorm:"not null"`
	PlanID    uuid.UUID      `gorm:"type:uuid;not null"`
	Plan      plans.Plan     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PlanID;references:ID;"`
	Password  string         `gorm:"not null"`
	Status    bool           `gorm:"not null;default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateUserCore(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
