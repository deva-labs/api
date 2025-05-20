package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/creack/pty"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math"
	"net/http"
	"os"
	"os/exec"
	"skypipe/src/config"
	"skypipe/src/lib/interfaces"
	"skypipe/src/modules/key_token/models"
	users "skypipe/src/modules/users/models"
	verifications "skypipe/src/modules/verifications/models"
	"strconv"
	"strings"
	"time"
)

func Paginate(total int64, page, perPage int) (map[string]interface{}, error) {
	if perPage <= 0 {
		perPage = 1
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	var nextPage, prevPage *int
	if page < totalPages {
		next := page + 1
		nextPage = &next
	}
	if page > 1 {
		prev := page - 1
		prevPage = &prev
	}

	return map[string]interface{}{
		"current_page":   page,
		"items_per_page": perPage,
		"next_page":      nextPage,
		"previous_page":  prevPage,
		"total_count":    total,
		"total_pages":    totalPages,
	}, nil
}

func ConvertStringToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func ConvertInt64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func ConvertInt64ToUint(i int64) uint {
	return uint(i)
}

func ConvertStringToUint(str string) uint {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(i)
}

func ConvertStringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return i
}

type ServiceError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *ServiceError) Error() string {
	return e.Message
}

type CalculateOffsetStruct struct {
	CurrentPage  int
	ItemsPerPage int
	OrderBy      string
	SortBy       string
	Offset       int
}

func CalculateOffset(currentPage, itemsPerPage int, sortBy, orderBy string) CalculateOffsetStruct {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if sortBy != "asc" && sortBy != "desc" {
		sortBy = "desc"
	}

	offset := (currentPage - 1) * itemsPerPage
	if offset < 0 {
		offset = 0
	}

	return CalculateOffsetStruct{
		CurrentPage:  currentPage,
		ItemsPerPage: itemsPerPage,
		OrderBy:      orderBy,
		SortBy:       sortBy,
		Offset:       offset,
	}
}

// BindJson for Fiber
func BindJson(c *fiber.Ctx, request interface{}) *ServiceError {
	if err := c.BodyParser(request); err != nil {
		return &ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input",
			Err:        err,
		}
	}
	return nil
}

func ExecWithAnimation(sc *interfaces.SafeConn, msg, command, action string, envVars map[string]string) error {
	spinnerFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerIndex := 0
	startTime := time.Now()

	cmd := exec.Command("bash", "-c", command)
	cmd.Env = os.Environ()
	for key, value := range envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	ptmx, err := pty.Start(cmd)
	if err != nil {
		sc.SafeWrite(websocket.TextMessage, []byte(fmt.Sprintf("✗ Failed to start command: %v", err)))
		return fmt.Errorf("failed to start command '%s': %w", command, err)
	}
	defer ptmx.Close()

	done := make(chan error, 1)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				sc.SafeWrite(websocket.TextMessage, buf[:n])
			}
			if err != nil {
				break
			}
		}
		done <- cmd.Wait()
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var execErr error

loop:
	for {
		select {
		case <-ticker.C:
			spinner := spinnerFrames[spinnerIndex%len(spinnerFrames)]
			elapsed := time.Since(startTime).Seconds()
			line := fmt.Sprintf("\r%s %s %s... %.2fs      ", spinner, strings.Title(action), msg, elapsed)
			sc.SafeWrite(websocket.TextMessage, []byte(line))
			spinnerIndex++
		case err := <-done:
			elapsed := time.Since(startTime).Seconds()
			clearLine := fmt.Sprintf("\r%s\r", strings.Repeat(" ", 80))
			sc.SafeWrite(websocket.TextMessage, []byte(clearLine))

			if err == nil {
				actionDone := transformAction(action)
				successLine := fmt.Sprintf("✔ Successfully %s %s", actionDone, msg)
				sc.SafeWrite(websocket.TextMessage, []byte(successLine))
				sc.SafeWrite(websocket.TextMessage, []byte(fmt.Sprintf("\n⏱️ Step Completed in %.2fs\n", elapsed)))
			} else {
				failLine := fmt.Sprintf("✗ Failed to %s %s (after %.2fs)", action, msg, elapsed)
				sc.SafeWrite(websocket.TextMessage, []byte(failLine))
				sc.SafeWrite(websocket.TextMessage, []byte(fmt.Sprintf("Error details: %v", err)))
				execErr = fmt.Errorf("command failed: %w", err)
			}
			break loop
		}
	}

	return execErr
}

func transformAction(action string) string {
	if strings.HasSuffix(action, "ing") {
		return strings.TrimSuffix(action, "ing") + "ed"
	}
	return action
}

func HashToken(token string) string {
	hashedToken := sha256.Sum256([]byte(token))
	tokenKey := hex.EncodeToString(hashedToken[:])
	return tokenKey
}

func CheckPasswordHash(pw, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil
}

func VerifyRedisTokenWithUserID(userID, token string) (bool, error) {
	rdb := config.RDB
	ctx := config.Ctx
	clientToken := HashToken(token)
	storedHash, err := rdb.Get(ctx, userID).Result()
	if err != nil || clientToken != storedHash {
		return false, err
	}
	return true, nil
}

// VerifyCaptcha Stub example – replace with real API call to hCaptcha or Google reCAPTCHA
func VerifyCaptcha(token string) bool {
	return token == "pass" // or call external verifications
}

func CleanUserOldTokens(db *gorm.DB, email string, userID uuid.UUID) error {
	rootID, err := GetRootUserID(db)
	if err != nil {
		return fmt.Errorf("failed to get root user ID: %w", err)
	}

	// Set UpdatedBy manually via Model Update
	if err := db.Where("email = ?", email).Delete(&verifications.VerificationCode{}).Error; err != nil {
		return err
	}
	if err := db.Model(&key_token.AccessToken{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{"updated_by": rootID}).Delete(&key_token.AccessToken{}).Error; err != nil {
		return err
	}
	if err := db.Model(&key_token.RefreshToken{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{"updated_by": rootID}).Delete(&key_token.RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}

func GetRootUserID(db *gorm.DB) (uuid.UUID, error) {
	var user users.User
	if err := db.Where("name = ?", "root").First(&user).Error; err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}
