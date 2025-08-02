package captcha

import (
	"deva/src/config"
	"deva/src/functions"
	"deva/src/lib/dto"
	captcha "deva/src/modules/captcha/models"
	"deva/src/utils"
	"net/http"
)

func CreateCaptcha(body dto.CreateCaptchaRequest) (string, *utils.ServiceError) {
	db := config.DB

	var userCaptcha captcha.Captcha
	if err := db.Where("user_id = ?", body.UserID).Delete(&userCaptcha).Error; err != nil {
		return "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete captcha belongs to user",
			Err:        err,
		}
	}
	// Generate text
	captchaValue := functions.GenerateRandomCaptcha(6)

	// Generate image from text
	imgBase64, err := functions.GenerateCaptchaImageBase64(captchaValue)
	if err != nil {
		return "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to generate captcha image",
			Err:        err,
		}
	}

	// Save to DB
	newCaptcha := captcha.Captcha{
		UserID: body.UserID,
		Text:   captchaValue,
	}
	if err := db.Create(&newCaptcha).Error; err != nil {
		return "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to store captcha",
			Err:        err,
		}
	}

	return imgBase64, nil
}
