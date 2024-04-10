package handlers

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"calendar_scheduler/src/services"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func (h Handler) SendPasswordRecover(c *fiber.Ctx) error {
	email := c.Query(constants.Email)

	userRepository := repositories.NewUserRepository(h.db)
	user, _, err := userRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.UserNotFounded) {
			return NotFoundError(c)
		}
		return InternalServerError(c)
	}

	recoveryRepository := repositories.NewResetPasswordRepository(h.db)
	resetPassword := models.ResetPassword{
		UserId: user.Id,
		Email:  user.Email,
		Code:   uuid.NewString(),
		Expiry: time.Now().Add(time.Hour * 24),
	}
	if err = recoveryRepository.SetResetPassword(resetPassword); err != nil {
		return InternalServerError(c)
	}

	stmpService := services.NewSESService()
	emailData := services.EmailData{
		Email:   email,
		BaseUrl: c.BaseURL(),
		Code:    resetPassword.Code,
	}
	if err := stmpService.SendEmail(emailData); err != nil {
		return InternalServerError(c)
	}
	return ResponseOK(c,
		map[string]string{
			"message": "Successful!",
			"url":     fmt.Sprintf("http://localhost:3000/recover-password?code=%s", resetPassword.Code),
		},
	)
}
