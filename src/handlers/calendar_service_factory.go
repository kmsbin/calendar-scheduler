package handlers

import (
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"context"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendarService(c *fiber.Ctx) (*calendar.Service, *models.MessageHTTP) {
	token := c.Locals("token").(string)
	userId := c.Locals("user_id")
	if userId == nil {
		return nil, models.MessageHTTPFromFiberError(fiber.ErrUnauthorized)
	}
	user, err := repositories.NewUserRepository().GetDataFromToken(userId.(float64))
	if err != nil {
		return nil, &models.MessageHTTP{Message: "User not founded!", HttpCode: fiber.StatusUnauthorized}
	}
	config, err := repositories.GetGoogleAuthConfig()
	if err != nil {
		return nil, &models.MessageHTTP{Message: err.Error(), HttpCode: fiber.StatusInternalServerError}
	}
	client, err := getClient(token, user.Id, config)

	if err != nil {
		if tokenNotFoundedErr, ok := err.(calendarTokenNotFounded); ok {
			return nil, &models.MessageHTTP{Message: tokenNotFoundedErr.AuthUrl, HttpCode: fiber.StatusPreconditionRequired}
		}
		return nil, models.MessageHTTPFromFiberError(fiber.ErrInternalServerError)
	}
	service, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, &models.MessageHTTP{Message: err.Error(), HttpCode: fiber.StatusInternalServerError}
	}
	return service, nil
}
