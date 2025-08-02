package users

import (
	"deva/src/config"
	"deva/src/functions"
	"deva/src/lib/dto"
	key_token "deva/src/modules/key_token/services"
	plans "deva/src/modules/plans/models"
	roles "deva/src/modules/roles/models"
	users "deva/src/modules/users/models"
	verificationService "deva/src/modules/verifications/services"
	"deva/src/services"
	"deva/src/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

func GetUserInfo(userID uuid.UUID) (*users.User, *utils.ServiceError) {
	db := config.DB

	var user users.User

	// Preload profile based on defined relationship
	if err := db.Preload("Profile").Preload("Plan").First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utils.ServiceError{
				StatusCode: http.StatusBadRequest,
				Message:    "User not found",
				Err:        err,
			}
		}
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "DB error",
			Err:        err,
		}
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
	hashed, err := services.HashPassword(body.Password)
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
		"user_id": user.ID,
	}, nil
}

func LoginService(body dto.LoginRequest) (map[string]interface{}, *utils.ServiceError) {
	db := config.DB
	rdb := config.RDB
	ctx := config.Ctx

	// Step 1: Find user
	var user users.User
	if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid email or password",
			Err:        err,
		}
	}

	// Step 2: Check password
	if !utils.CheckPasswordHash(body.Password, user.Password) {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid email or password",
		}
	}

	// Step 3: Clean up old records (accessToken, refreshToken, verificationCode)
	if err := utils.CleanUserOldTokens(db, user.Email, user.ID); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to clear user session",
			Err:        err,
		}
	}

	// Step 4: Generate and save verification code
	verificationCode := verificationService.GenerateVerificationCode()
	if err := verificationService.SaveVerificationCode(user.Email, verificationCode); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save verification code",
			Err:        err,
		}
	}

	// Step 5: Generate and store secure token (Redis)
	token, err := key_token.GenerateSecureToken(32)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to generate token",
			Err:        err,
		}
	}
	redisKey := fmt.Sprintf("login_token:%s", user.ID.String())
	tokenKey := utils.HashToken(token)

	ttl := time.Minute * 5
	if ttlEnv := os.Getenv("RESET_TOKEN_TTL"); ttlEnv != "" {
		if parsed, err := time.ParseDuration(ttlEnv); err == nil {
			ttl = parsed
		}
	}

	if err := rdb.Set(ctx, redisKey, tokenKey, ttl).Err(); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Could not store token",
			Err:        err,
		}
	}

	// Step 6: Send email
	if err := verificationService.SendVerificationEmail(user.Email, verificationCode); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to send verification email",
			Err:        err,
		}
	}

	// Return login metadata
	return map[string]interface{}{
		"token":   token,
		"user_id": user.ID,
		"email":   user.Email,
	}, nil
}

// ForgotPassword is the function will receive an email to process forgot password services
func ForgotPassword(email string) (map[string]interface{}, *utils.ServiceError) {
	db := config.DB
	rdb := config.RDB
	ctx := config.Ctx
	// Step 1: Check user?
	var user users.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid email or password",
			Err:        err,
		}
	}
	// Step 2: Clean up old records (accessToken, refreshToken, verificationCode)
	if err := utils.CleanUserOldTokens(db, user.Email, user.ID); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to clear user session",
			Err:        err,
		}
	}

	// Step 3: Generate and save verification code
	verificationCode := verificationService.GenerateVerificationCode()
	if err := verificationService.SaveVerificationCode(user.Email, verificationCode); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save verification code",
			Err:        err,
		}
	}

	// Step 4: Generate and store secure token (Redis)
	token, err := key_token.GenerateSecureToken(32)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to generate token",
			Err:        err,
		}
	}
	tokenKey := utils.HashToken(token)

	ttl := time.Minute * 5
	if ttlEnv := os.Getenv("RESET_TOKEN_TTL"); ttlEnv != "" {
		if parsed, err := time.ParseDuration(ttlEnv); err == nil {
			ttl = parsed
		}
	}
	redisKey := fmt.Sprintf("reset_password_token:%s", user.ID.String())
	if err := rdb.Set(ctx, redisKey, tokenKey, ttl).Err(); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Could not store token",
			Err:        err,
		}
	}

	// Step 5: Send email
	if err := verificationService.SendVerificationEmail(user.Email, verificationCode); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to send verification email",
			Err:        err,
		}
	}

	// Return login metadata
	return map[string]interface{}{
		"token":   token,
		"user_id": user.ID,
		"email":   user.Email,
	}, nil
}

// RenewPassword is the function to change the password follow by user
func RenewPassword(request dto.RenewPasswordRequest) *utils.ServiceError {
	// Start a database transaction
	tx := config.DB.Begin()
	var user users.User

	// Step 1: Check if the user exists
	if err := tx.Where("id = ?", request.UserID).First(&user).Error; err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "User not found",
			Err:        err,
		}
	}
	// Step 2: Validate with redis db
	redisKey := fmt.Sprintf("reset_password_token:%s", user.ID.String())
	validated, err := utils.VerifyRedisTokenWithUserID(redisKey, request.Token)
	if err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Token validation error",
			Err:        err,
		}
	}
	if !validated {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid or expired token",
			Err:        errors.New("token is not valid"),
		}
	}

	// Step 3: Validate Captcha
	if err := functions.ValidateCaptcha(request.Captcha, request.UserID); err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid captcha",
			Err:        err,
		}
	}
	// Step 3: Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 12)
	if err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to hash new password",
			Err:        err,
		}
	}

	// Step 4: Update the password in the database
	user.Password = string(hashedPassword)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save new password",
			Err:        err,
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to commit transaction",
			Err:        err,
		}
	}

	// Step 5: Return success response
	return nil
}

// ChangePassword is the function to change the password follow by user
func ChangePassword(oldPassword, newPassword string, userID uuid.UUID) *utils.ServiceError {
	// Start a database transaction
	tx := config.DB.Begin()
	var user users.User

	// Step 1: Check if the user exists
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "User not found",
			Err:        err,
		}
	}

	// Step 2: Validate the old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid old password",
			Err:        err,
		}
	}

	// Step 3: Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to hash new password",
			Err:        err,
		}
	}

	// Step 4: Update the password in the database
	user.Password = string(hashedPassword)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save new password",
			Err:        err,
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to commit transaction",
			Err:        err,
		}
	}

	// Step 5: Return success response
	return nil
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
