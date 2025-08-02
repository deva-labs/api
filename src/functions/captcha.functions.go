package functions

import (
	"bytes"
	"deva/src/config"
	captcha "deva/src/modules/captcha/models"
	"encoding/base64"
	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"image/color"
	"image/png"
	"math/rand"
)

func GenerateRandomCaptcha(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func GenerateCaptchaImageBase64(text string) (string, error) {
	const width = 200
	const height = 70

	dc := gg.NewContext(width, height)

	// Background
	dc.SetRGB(1, 1, 1) // white
	dc.Clear()

	// Random lines (noise)
	for i := 0; i < 8; i++ {
		dc.SetRGBA(rand.Float64(), rand.Float64(), rand.Float64(), 0.5)
		x1 := rand.Float64() * width
		y1 := rand.Float64() * height
		x2 := rand.Float64() * width
		y2 := rand.Float64() * height
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	// Draw captcha text
	if err := dc.LoadFontFace("/usr/share/fonts/dejavu/DejaVuSans-Bold.ttf", 36); err != nil {
		return "", err
	}
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(text, width/2, height/2, 0.5, 0.5)

	// Encode image to base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return "", err
	}

	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Img, nil
}

func ValidateCaptcha(text string, userID uuid.UUID) error {
	db := config.DB
	var userCaptcha captcha.Captcha
	if err := db.Where("user_id = ? and text = ?", userID, text).First(&userCaptcha).Error; err != nil {
		return err
	}
	return nil
}
