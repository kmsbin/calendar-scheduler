package handlers

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"github.com/gofiber/fiber/v2"
)

func (h Handler) DeleteUser(c *fiber.Ctx) error {
	userId := c.Locals(constants.UserId)
	userRepository := repositories.NewUserRepository(h.db)
	err := userRepository.DeleteUser(userId)

	if err != nil {
		return models.MessageHTTPFromFiberError(fiber.ErrInternalServerError).FiberContext(c)
	}

	return c.
		Status(fiber.StatusOK).
		JSON(models.MessageHTTPFromMessage("Success"))
}
