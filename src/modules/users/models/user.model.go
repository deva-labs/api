package users

import (
	plans "deva/src/modules/plans/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Name      string         `gorm:"not null"`
	PlanID    uuid.UUID      `gorm:"type:uuid;not null"`
	Plan      plans.Plan     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PlanID;references:ID;"`
	ProfileID uuid.UUID      `gorm:"type:uuid"`
	Profile   *Profile       `gorm:"foreignKey:ProfileID;references:ID"`
	Password  string         `gorm:"not null"`
	Status    bool           `gorm:"not null;default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateUserCore(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
