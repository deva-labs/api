package users

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"skypipe/src/config"
	"skypipe/src/lib/dto"
	plans "skypipe/src/modules/plans/models"
	roles "skypipe/src/modules/roles/models"
	users "skypipe/src/modules/users/models"
	verificationService "skypipe/src/modules/verifications/services"
	"skypipe/src/utils"
)

func GetUserInfo(userID uuid.UUID) (*users.User, error) {
	db := config.DB

	var user users.User

	// Preload profile based on defined relationship
	if err := db.Preload("Profile").First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func RegisterService(body dto.RegisterRequest) (map[string]interface{}, *utils.ServiceError) {
	db := config.DB

	// Check if user already exists
	var existing users.User
	if err := db.Where("email = ?", body.Email).First(&existing).Error; err == nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusConflict,
			Message:    "Email already registered",
			Err:        err,
		}
	}

	// Hash password
	hashed, err := utils.HashPassword(body.Password)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Hashing failed",
			Err:        err,
		}
	}

	// Check plan
	var selectedPlan plans.Plan
	if err := db.First(&selectedPlan, "id = ?", body.PlanID).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid plan ID",
			Err:        err,
		}
	}

	// Create new user
	user := users.User{
		Name:     body.Name,
		Email:    body.Email,
		PlanID:   selectedPlan.ID,
		Password: hashed,
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Registration failed",
			Err:        err,
		}
	}

	// Create role for new user
	var roleUser roles.Role
	if err := db.Where("name = ?", "user").First(&roleUser).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Role select failed",
			Err:        err,
		}
	}
	userRole := users.UserRole{
		UserID:    user.ID,
		RoleID:    roleUser.ID,
		UpdatedBy: user.ID,
	}
	if err := db.Create(&userRole).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "User role create failed",
			Err:        err,
		}
	}
	// Create profile
	profile := users.Profile{
		UserID:    user.ID,
		FullName:  body.FullName,
		Phone:     body.Phone,
		Gender:    body.Gender,
		Country:   body.Country,
		City:      body.City,
		UpdatedBy: user.ID,
	}
	if err := db.Create(&profile).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create user profile",
			Err:        err,
		}
	}

	return map[string]interface{}{
		"message": "User registered successfully",
		"user_id": user.ID,
	}, nil
}

func LoginService(body dto.LoginRequest) (map[string]interface{}, *utils.ServiceError) {
	db := config.DB
	// Optionally use rate limiting middleware before this handler

	var user users.User
	if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid email or password",
			Err:        err,
		}
	}

	if !utils.CheckPasswordHash(body.Password, user.Password) {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid email or password",
		}
	}

	// Send verifications code
	verificationCode := verificationService.GenerateVerificationCode()

	// Save new verifications code to db
	if err := verificationService.SaveVerificationCode(user.Email, verificationCode); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save verifications code",
			Err:        err,
		}
	}
	if err := verificationService.SendVerificationEmail(user.Email, verificationCode); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to send verifications code",
			Err:        err,
		}
	}
	return map[string]interface{}{
		"email": user.Email,
	}, nil
}

// ForgotPassword is the function will receive an email to process forgot password service
func ForgotPassword(email string) (int, error) {
	db := config.DB
	// Step 1: Check user?
	var user users.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return http.StatusBadRequest, errors.New("user not found")
	}
	// Step 2: Generate 6 digits code
	verificationCode := verificationService.GenerateVerificationCode()

	// Step 3: Save the verifications code to the database
	if err := verificationService.SaveVerificationCode(email, verificationCode); err != nil {
		return http.StatusInternalServerError, errors.New("failed to save verifications code")
	}

	// Step 4: Send the verifications code via email
	if err := verificationService.SendVerificationEmail(email, verificationCode); err != nil {
		return http.StatusInternalServerError, errors.New("failed to send verifications email")
	}

	return http.StatusOK, nil
}

// RenewPassword is the function to change the password follow by user
func RenewPassword(newPassword, userID string) (int, error) {
	// Start a database transaction
	tx := config.DB.Begin()
	var user users.User

	// Step 1: Check if the user exists
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		return http.StatusNotFound, errors.New("user not found")
	}

	// Step 2: Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to hash new password")
	}

	// Step 4: Update the password in the database
	user.Password = string(hashedPassword)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to save new password")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to commit transaction")
	}

	// Step 5: Return success response
	return http.StatusOK, nil
}

// ChangePassword is the function to change the password follow by user
func ChangePassword(oldPassword, newPassword, userID string) (int, error) {
	// Start a database transaction
	tx := config.DB.Begin()
	var user users.User

	// Step 1: Check if the user exists
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		return http.StatusNotFound, errors.New("user not found")
	}

	// Step 2: Validate the old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		tx.Rollback()
		return http.StatusUnauthorized, errors.New("old password is incorrect")
	}

	// Step 3: Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to hash new password")
	}

	// Step 4: Update the password in the database
	user.Password = string(hashedPassword)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to save new password")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to commit transaction")
	}

	// Step 5: Return success response
	return http.StatusOK, nil
}

func GetUserImageByID(userID uint) (string, *utils.ServiceError) {
	if userID == 0 {
		return "", &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid user ID",
		}
	}

	var userAvatar string
	err := config.DB.
		Model(&users.Profile{}).
		Select("avatar_url").
		Where("user_id = ?", userID).
		Scan(&userAvatar).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", &utils.ServiceError{
				StatusCode: http.StatusNotFound,
				Message:    "user profile not found",
			}
		}

		return "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to get user image: %v", err),
		}
	}

	// Return empty string if no profile picture is set
	return userAvatar, nil
}
