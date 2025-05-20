package key_token

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"skypipe/src/config"
	key_token "skypipe/src/modules/key_token/models"
	users "skypipe/src/modules/users/models"
	"skypipe/src/utils"
)

func GenerateHexTokens(userID uuid.UUID) (string, string, error) {
	db := config.DB

	accessTokenBytes := make([]byte, 32)
	refreshTokenBytes := make([]byte, 32)

	// Get root user
	rootUid, err := utils.GetRootUserID(db)
	if err != nil {
		return "", "", err
	}
	if _, err := rand.Read(accessTokenBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	accessToken := hex.EncodeToString(accessTokenBytes)
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	if err := db.Create(&key_token.AccessToken{
		UserID:    userID,
		Token:     accessToken,
		UpdatedBy: rootUid,
	}).Error; err != nil {
		return "", "", fmt.Errorf("failed to save access token: %w", err)
	}
	if err := db.Create(&key_token.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		UpdatedBy: rootUid,
	}).Error; err != nil {
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func DeleteAllTokensByUserID(userID uuid.UUID) *utils.ServiceError {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to begin transaction",
			Err:        tx.Error,
		}
	}

	if err := tx.Where("user_id = ?", userID).Delete(&key_token.AccessToken{}).Error; err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to delete access tokens",
			Err:        err,
		}
	}
	if err := tx.Where("user_id = ?", userID).Delete(&key_token.RefreshToken{}).Error; err != nil {
		tx.Rollback()
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to delete refresh tokens",
			Err:        err,
		}
	}
	if err := tx.Commit().Error; err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to commit transaction",
			Err:        err,
		}
	}
	return nil
}

func RefreshAccessToken(token, clientID string) (map[string]interface{}, *utils.ServiceError) {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to begin transaction",
			Err:        tx.Error,
		}
	}

	var refreshToken key_token.RefreshToken
	if err := tx.Where("token = ? AND user_id = ?", token, clientID).First(&refreshToken).Error; err != nil {
		tx.Rollback()
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    "invalid or expired refresh token",
			Err:        err,
		}
	}

	if err := tx.Where("user_id = ?", refreshToken.UserID).Delete(&key_token.AccessToken{}).Error; err != nil {
		tx.Rollback()
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to delete old access tokens",
			Err:        err,
		}
	}

	newToken, err := GenerateAccessToken()
	if err != nil {
		tx.Rollback()
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to generate new access token",
			Err:        err,
		}
	}

	newAccessToken := key_token.AccessToken{
		Token:  newToken,
		UserID: refreshToken.UserID,
		Status: true,
	}
	if err := tx.Create(&newAccessToken).Error; err != nil {
		tx.Rollback()
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to create new access token",
			Err:        err,
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to commit transaction",
			Err:        err,
		}
	}

	response := map[string]interface{}{
		"access_token": newAccessToken.Token,
		"expires_at":   newAccessToken.ExpiresAt,
	}
	return response, nil
}

func GenerateAccessToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return hex.EncodeToString(tokenBytes), nil
}

func VerifyToken(token string) (*users.User, *utils.ServiceError) {
	db := config.DB
	var tokenRecord key_token.AccessToken
	if err := db.Where("token = ?", token).First(&tokenRecord).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Message:    "token not found",
			Err:        err,
		}
	}

	var user users.User
	if err := db.Where("id = ?", tokenRecord.UserID).First(&user).Error; err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusNotFound,
			Message:    "user not found",
			Err:        err,
		}
	}

	return &user, nil
}

func GenerateSecureToken(byteLength int) (string, error) {
	if byteLength < 16 {
		return "", fmt.Errorf("token length too short; must be >= 16 bytes")
	}

	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
