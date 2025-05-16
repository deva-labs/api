package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Name     string    `json:"name" validate:"required,min=2"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8"`
	PlanID   uuid.UUID `json:"plan_id" validate:"required,uuid4"`
	Captcha  string    `json:"captcha"`

	// Minimum profile fields
	FullName string `json:"full_name" validate:"required,min=2"`
	Phone    string `json:"phone" validate:"required,min=10"`
	Gender   string `json:"gender" validate:"omitempty,oneof=Male Female Other"`
	Country  string `json:"country" validate:"omitempty"`
	City     string `json:"city" validate:"omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PermissionInterface struct {
	Resource string
	Action   string
}
