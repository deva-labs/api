package functions

import (
	"deva/src/lib/dto"
	users "deva/src/modules/users/models"
)

func ToUserResponse(u *users.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Plan: dto.PlanResponse{
			ID:                   u.Plan.ID,
			Name:                 u.Plan.Name,
			Price:                u.Plan.Price,
			MaxProjects:          u.Plan.MaxProjects,
			MaxDeploymentsPerDay: u.Plan.MaxDeploymentsPerDay,
		},
		Profile: dto.ProfileResponse{
			ID:        u.Profile.ID,
			FullName:  u.Profile.FullName,
			Bio:       u.Profile.Bio,
			AvatarURL: u.Profile.AvatarURL,
			Gender:    u.Profile.Gender,
			Phone:     u.Profile.Phone,
			Address:   u.Profile.Address,
			Country:   u.Profile.Country,
			City:      u.Profile.City,
		},
	}
}
