package dto

import (
	"github.com/google/uuid"
	"time"
)

type UserResponse struct {
	ID        uuid.UUID       `json:"id"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Status    bool            `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Plan      PlanResponse    `json:"plan"`
	Profile   ProfileResponse `json:"profile"`
}

type PlanResponse struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	Price                float64   `json:"price"`
	MaxProjects          int       `json:"max_projects"`
	MaxDeploymentsPerDay int       `json:"max_deployments_per_day"`
}

type ProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	Gender    string    `json:"gender"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
}
