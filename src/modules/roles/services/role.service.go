package roles

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"skypipe/src/config"
	roles "skypipe/src/modules/roles/models"
)

// GetRoleByUserID retrieves the role of a user by their user ID
func GetRoleByUserID(userID uuid.UUID) (*roles.Role, error) {
	db := config.DB

	var role roles.Role
	err := db.
		Table("roles").
		Select("roles.*").
		Joins("inner join user_roles on user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		First(&role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found for this user")
		}
		return nil, err
	}

	return &role, nil
}

// CreateRoleByAdmin creates a new role with the given name
func CreateRoleByAdmin(roleName string) (fiber.Map, error) {
	db := config.DB

	if roleName == "" {
		return nil, errors.New("role name cannot be empty")
	}

	var existing roles.Role
	if err := db.Where("name = ?", roleName).First(&existing).Error; err == nil {
		return nil, errors.New("role already exists")
	}

	role := &roles.Role{Name: roleName}
	if err := db.Create(role).Error; err != nil {
		return nil, err
	}

	return fiber.Map{
		"id":   role.ID,
		"name": role.Name,
	}, nil
}

// ListAllRoles returns all roles in the system
func ListAllRoles() ([]fiber.Map, error) {
	db := config.DB

	var roles []roles.Role
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}

	var result []fiber.Map
	for _, r := range roles {
		result = append(result, fiber.Map{
			"id":   r.ID,
			"name": r.Name,
		})
	}
	return result, nil
}
