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
	"log"
	"time"
)

func (h Handler) SendPasswordRecover(c *fiber.Ctx) error {
	email := c.Query(constants.Email)

	userRepository := repositories.NewUserRepository(h.db)
	user, _, err := userRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.UserNotFounded) {
			return c.
				Status(fiber.StatusNotFound).
				JSON(models.MessageHTTPFromFiberError(fiber.ErrNotFound))
		}
		return models.
			MessageHTTPFromFiberError(fiber.ErrInternalServerError).
			FiberContext(c)
	}

	recoveryRepository := repositories.NewResetPasswordRepository(h.db)
	resetPassword := models.ResetPassword{
		UserId: user.Id,
		Email:  user.Email,
		Code:   uuid.NewString(),
		Expiry: time.Now().Add(time.Hour * 24),
	}
	err = recoveryRepository.SetResetPassword(resetPassword)
	if err != nil {
		return models.MessageHTTPFromFiberError(fiber.ErrInternalServerError).FiberContext(c)
	}

	log.Println(user)
	stmpService := services.NewSESService()
	if err := stmpService.SendEmail(email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
	}

	return c.
		Status(200).
		JSON(map[string]string{
			"message": "Successful!",
			"url":     fmt.Sprintf("http://localhost:3000/recover-password?code=%s", resetPassword.Code),
		})
}
