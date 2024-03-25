package handlers

import (
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"errors"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) ReceivePasswordRecover(c *fiber.Ctx) error {
	data := struct {
		Password string `json:"password"`
		Code     string `json:"code"`
	}{}

	err := c.BodyParser(&data)

	if data.Code == "" || data.Password == "" {
		return c.
			Status(fiber.StatusUnprocessableEntity).
			JSON(models.MessageHTTPFromMessage("missing parameter"))
	}

	resetRepository := repositories.NewResetPasswordRepository(h.db)
	resetPasswordData, err := resetRepository.GetResetPasswordByCode(data.Code)

	if err != nil {
		if errors.Is(repositories.ResetPasswordNotFound, err) {
			return c.
				Status(fiber.StatusNotFound).
				JSON(models.MessageHTTPFromMessage(repositories.ResetPasswordNotFound.Error()))
		}
		return InternalServerError(c)
	}
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if err != nil {
		return InternalServerError(c)
	}
	userRepository := repositories.NewUserRepository(h.db)
	err = userRepository.ResetPassword(resetPasswordData.UserId, string(passwordBytes))
	if err != nil {
		return InternalServerError(c)
	}

	err = resetRepository.DeleteResetPasswordData(resetPasswordData.Code)

	if err != nil {
		return InternalServerError(c)
	}

	return c.JSON(models.MessageHTTPFromMessage("Password reseted sucessful!"))
}
