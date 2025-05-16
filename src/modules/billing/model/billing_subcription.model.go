package billing

import (
	users "skypipe/src/modules/users/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingSubscription struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null"`
	User           users.User     `gorm:"foreignKey:UserID;references:ID"`
	Provider       string         `gorm:"not null"`
	SubscriptionID string         `gorm:"not null"`
	Status         string         `gorm:"not null"`
	StartedAt      time.Time      `gorm:"not null"`
	RenewAt        time.Time      `gorm:"not null"`
	UpdatedBy      uuid.UUID      `gorm:"type:uuid;not null"`
	UpdatedByUser  users.User     `gorm:"foreignKey:UpdatedBy;references:ID"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func MigrateBillingSubscriptions(db *gorm.DB) error {
	return db.AutoMigrate(&BillingSubscription{})
}
