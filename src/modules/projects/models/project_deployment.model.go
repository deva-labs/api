package projects

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	deployments "skypipe/src/modules/deployments/models"
	users "skypipe/src/modules/users/models"
	"time"
)

type ProjectDeployment struct {
	ID            uuid.UUID              `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID     uuid.UUID              `gorm:"type:uuid;not null"`
	Project       Project                `gorm:"foreignKey:ProjectID;references:ID"`
	DeploymentID  uuid.UUID              `gorm:"type:uuid;not null"`
	Deployment    deployments.Deployment `gorm:"foreignKey:DeploymentID;references:ID"`
	UpdatedBy     uuid.UUID              `gorm:"type:uuid;not null"`
	UpdatedByUser users.User             `gorm:"foreignKey:UpdatedBy;references:ID;"`
	CreatedAt     time.Time              `gorm:"autoCreateTime"`
	UpdatedAt     time.Time              `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt         `gorm:"index"`
}

func MigrateProjectDeployment(db *gorm.DB) error {
	return db.AutoMigrate(&ProjectDeployment{})
}
