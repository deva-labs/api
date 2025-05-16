package key_token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"skypipe/src/config"
	key_token "skypipe/src/modules/key_token/models"
	users "skypipe/src/modules/users/models"
)

func GenerateHexTokens(userID uuid.UUID) (string, string, error) {
	db := config.DB

	// Generate random hex strings for tokens
	accessTokenBytes := make([]byte, 32) // 256-bit
	refreshTokenBytes := make([]byte, 32)

	_, err := rand.Read(accessTokenBytes)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	accessToken := hex.EncodeToString(accessTokenBytes)
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	// Save access token to database
	accessTokenEntry := key_token.AccessToken{
		UserID: userID,
		Token:  accessToken,
	}
	if err := db.Create(&accessTokenEntry).Error; err != nil {
		return "", "", err
	}

	// Save refresh token to database
	refreshTokenEntry := key_token.RefreshToken{
		UserID: userID,
		Token:  refreshToken,
	}
	if err := db.Create(&refreshTokenEntry).Error; err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func DeleteAllTokensByUserID(userID uuid.UUID) error {
	// Start a transaction to ensure both deletions are atomic
	tx := config.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	// Delete all AccessTokens for the given UserID
	if err := tx.Where("user_id = ?", userID).Delete(&key_token.AccessToken{}).Error; err != nil {
		tx.Rollback() // Rollback the transaction if deletion fails
		return fmt.Errorf("failed to delete access tokens: %v", err)
	}

	// Delete all RefreshTokens for the given UserID
	if err := tx.Where("user_id = ?", userID).Delete(&key_token.RefreshToken{}).Error; err != nil {
		tx.Rollback() // Rollback the transaction if deletion fails
		return fmt.Errorf("failed to delete refresh tokens: %v", err)
	}

	// Commit the transaction if both deletions succeed
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// RefreshAccessToken is the function to renew the access token by refresh token
func RefreshAccessToken(token, clientID string) (map[string]interface{}, error) {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	// Check if the refresh token is valid
	var refreshToken key_token.RefreshToken
	if err := tx.Where("token = ? AND user_id = ?", token, clientID).First(&refreshToken).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("invalid or expired refresh token: %v", err)
	}

	// Delete old access tokens for the user
	if err := tx.Where("user_id = ?", refreshToken.UserID).Delete(&key_token.AccessToken{}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to delete old access tokens: %v", err)
	}

	// Generate a new access token
	token, err := GenerateAccessToken()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to generate new access token: %v", err)
	}

	newAccessToken := key_token.AccessToken{
		Token:  token,
		UserID: refreshToken.UserID,
		Status: true,
	}

	if err := tx.Create(&newAccessToken).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create new access token: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Return the new access token
	response := map[string]interface{}{
		"access_token": newAccessToken.Token,
		"expires_at":   newAccessToken.ExpiresAt,
	}

	return response, nil
}

// GenerateAccessToken is the function to generate an access token
func GenerateAccessToken() (string, error) {
	// Generate random bytes for the token
	tokenBytes := make([]byte, 32) // 256-bit token
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}

	// Encode the random bytes to a hexadecimal string
	accessToken := hex.EncodeToString(tokenBytes)

	return accessToken, nil
}

// VerifyToken is a function help to verify and return user from token
func VerifyToken(token string) (*users.User, error) {
	db := config.DB
	// Proceed with token validation
	var tokenRecord key_token.AccessToken
	if err := db.Where("token = ?", token).First(&tokenRecord).Error; err != nil {
		return nil, err
	}

	// Find User By ID
	var user users.User
	if err := db.Where("id = ?", tokenRecord.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
