package handlers

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"github.com/gofiber/fiber/v2"
	"log"
)

func (h Handler) SignOutUser(c *fiber.Ctx) error {
	if httpModel := ValidateToken(c.Query(constants.Token), c); httpModel != nil {
		return httpModel.FiberContext(c)
	}

	tokenData, httpError := GetTokenExpirationData(c.Query(constants.Token), c)
	if httpError != nil {
		return httpError.FiberContext(c)
	}
	authRepository := repositories.NewAuthRepository(h.db)

	err := authRepository.InsertTokenBlackList(
		tokenData.UserId,
		tokenData.Token,
		tokenData.Expiration,
	)
	if err != nil {
		log.Println(err)
		return models.MessageHTTPFromFiberError(fiber.ErrInternalServerError).FiberContext(c)
	}

	return c.
		Status(fiber.StatusOK).
		JSON(models.MessageHTTPFromMessage("sign out successful!"))
}
