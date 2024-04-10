package handlers

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/repositories"
	"github.com/gofiber/fiber/v2"
	"log"
)

func (h Handler) SignOutUser(c *fiber.Ctx) error {
	token := c.Query(constants.Token)

	if httpModel := ValidateToken(token, c); httpModel != nil {
		return httpModel.FiberContext(c)
	}

	tokenData, httpError := GetTokenExpirationData(token, c)
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
		return InternalServerError(c)
	}

	return ResponseOK(c)
}
