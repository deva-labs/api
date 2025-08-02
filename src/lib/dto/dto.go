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

type CreateCaptchaRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid4"`
}
type PermissionInterface struct {
	Resource string
	Action   string
}

type VerificationCodeRequest struct {
	Token  string    `json:"token" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required,uuid4"`
	Email  string    `json:"email" validate:"required,email"`
	Code   string    `json:"code" validate:"required"`
}

type RenewPasswordRequest struct {
	NewPassword string    `json:"new_password" validate:"required,min=8"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid4"`
	Token       string    `json:"token" validate:"required"`
	Captcha     string    `json:"captcha" validate:"required,min=6"`
}

type CreateFiberRequest struct {
	ProjectName string            `json:"project_name"`
	UserID      uuid.UUID         `json:"user_id"`
	Env         map[string]string `json:"env"`
}

type ChangePasswordRequest struct {
	OldPassword string    `json:"old_password" validate:"required,min=8"`
	NewPassword string    `json:"new_password" validate:"required,min=8"`
	UserID      uuid.UUID `json:"user_id"`
}
