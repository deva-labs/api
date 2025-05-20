package verifications

import (
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"skypipe/src/config"
	"skypipe/src/lib/dto"
	key_token "skypipe/src/modules/key_token/services"
	users "skypipe/src/modules/users/models"
	verifications "skypipe/src/modules/verifications/models"
	"skypipe/src/utils"
	"time"
)

// GenerateVerificationCode Helper function to generate a 6-digit verifications code
func GenerateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// SaveVerificationCode Helper function to save the verifications code
func SaveVerificationCode(email, code string) error {
	db := config.DB
	var user users.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	// Delete all verifications code before sent a new one
	if err := db.Where("email = ?", email).Delete(&verifications.VerificationCode{}).Error; err != nil {
		return fmt.Errorf("failed to delete verifications code: %v", err)
	}

	verification := verifications.VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}
	return db.Create(&verification).Error
}

// SendVerificationEmail sends a verifications email with the given code
func SendVerificationEmail(email, code string) error {
	// Email server configuration
	smtpHost := os.Getenv("SMTP_HOST")     // E.g., "smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT")     // E.g., "587"
	smtpUsername := os.Getenv("SMTP_USER") // Your email address
	smtpPassword := os.Getenv("SMTP_PASS") // Your email password or app-specific password

	// Email content
	subject := "Your Verification Code"
	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		email, subject, fmt.Sprintf("<h2>Your verification code is: %s</h2>", code)))

	// Set up authentication information
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		auth,
		smtpUsername,
		[]string{email},
		message,
	)
	if err != nil {
		return fmt.Errorf("failed to send verifications email: %w", err)
	}

	return nil
}

// VerifyCodeAndGenerateTokens Function to verify the code and generate tokens
func VerifyCodeAndGenerateTokens(code dto.VerificationCodeRequest) (map[string]interface{}, *utils.ServiceError) {
	db := config.DB
	maxOtpAttempts := utils.ConvertStringToInt(os.Getenv("MAX_OTP_ATTEMPTS"))

	// Step 1: Retrieve the user
	var user users.User
	if err := db.Where("id = ?", code.UserID).First(&user).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "User not found",
			Err:        err,
		}
	}

	// Step 2: Retrieve verification code
	var verification verifications.VerificationCode
	if err := db.Where("email = ?", code.Email).First(&verification).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Verification record not found",
			Err:        err,
		}
	}

	// Step 3: Check if expired
	if verification.ExpiresAt.Before(time.Now()) {
		db.Delete(&verification)
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Verification code expired",
			Err:        nil,
		}
	}

	// Step 4: Validate code
	if verification.Code != code.Code {
		verification.InputCount++
		db.Save(&verification)

		if verification.InputCount >= maxOtpAttempts {
			db.Delete(&verification)
			user.Status = false
			db.Save(&user)

			return nil, &utils.ServiceError{
				StatusCode: http.StatusBadRequest,
				Message:    "Too many failed attempts; your account has been suspended. Please contact support",
				Err:        nil,
			}
		}

		remaining := maxOtpAttempts - verification.InputCount
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("invalid code. You have %d attempt(s) remaining", remaining),
			Err:        nil,
		}
	}

	// Step 5: Cleanup verification record
	db.Delete(&verification)

	// Step 6: Invalidate existing tokens
	if err := key_token.DeleteAllTokensByUserID(user.ID); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "could not invalidate existing sessions",
			Err:        err,
		}
	}

	// Step 7: Generate new access/refresh tokens
	accessToken, refreshToken, err := key_token.GenerateHexTokens(user.ID)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "could not generate verification code",
		}
	}

	// Step 8: Return token response
	return map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func VerifyCodeAndSetPasswordToken(code, email string) (map[string]interface{}, *utils.ServiceError) {
	db := config.DB
	rdb := config.RDB
	ctx := config.Ctx
	maxOtpAttempts := utils.ConvertStringToInt(os.Getenv("MAX_OTP_ATTEMPTS"))
	var user users.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("User not found"),
			Err:        err,
		}
	}

	var verification verifications.VerificationCode
	if err := db.Where("email = ? AND code = ?", email, code).First(&verification).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Invalid verification code"),
			Err:        err,
		}
	}

	// Check expired
	if verification.ExpiresAt.Before(time.Now()) {
		db.Delete(&verification)
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    fmt.Sprintf("Invalid verification code"),
			Err:        nil,
		}
	}

	// Check code match
	if verification.Code != code {
		// Use transaction to avoid race condition
		err := db.Transaction(func(tx *gorm.DB) error {
			verification.InputCount += 1
			if err := tx.Save(&verification).Error; err != nil {
				return err
			}

			if verification.InputCount >= maxOtpAttempts {
				if err := tx.Delete(&verification).Error; err != nil {
					return err
				}
				user.Status = false
				if err := tx.Save(&user).Error; err != nil {
					return err
				}
				return fmt.Errorf("too many attempts; account suspended")
			}
			return fmt.Errorf("invalid verification code")
		})
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    fmt.Sprintf("Invalid verification code"),
			Err:        err,
		}
	}

	// Delete old tokens
	if err := key_token.DeleteAllTokensByUserID(user.ID); err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("could not invalidate existing sessions"),
			Err:        err,
		}
	}

	// Generate new token
	token, err := key_token.GenerateSecureToken(32)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("could not generate authentication tokens"),
			Err:        err,
		}
	}

	// Hash token before storing to Redis
	tokenKey := utils.HashToken(token)

	// TTL from env (default 5m)
	ttl := time.Minute * 5
	if ttlEnv := os.Getenv("RESET_TOKEN_TTL"); ttlEnv != "" {
		if parsed, err := time.ParseDuration(ttlEnv); err == nil {
			ttl = parsed
		}
	}

	// Store token in Redis
	err = rdb.Set(ctx, user.ID.String(), tokenKey, ttl).Err()
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("could not save token"),
			Err:        err,
		}
	}

	return map[string]interface{}{
		"token":   token,
		"user_id": user.ID,
	}, nil
}
