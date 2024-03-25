package handlers

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"calendar_scheduler/src/services"
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func (h Handler) SetTokenGoogleCalendar(c *fiber.Ctx) error {
	token := c.Query("state")
	code := c.Query("code")
	if len(token) == 0 || len(code) == 0 {
		log.Printf("state %s", token)
		return fiber.ErrUnauthorized
	}

	httpModel := ValidateToken(token, c)
	if httpModel != nil {
		return httpModel.FiberContext(c)
	}
	userId := c.Locals(constants.UserId)

	if userId == nil {
		return fiber.ErrUnauthorized
	}
	userRepository := repositories.NewUserRepository(h.db)
	_, err := userRepository.GetUserById(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.MessageHTTP{Message: "User not founded!"})
	}
	config, err := repositories.GetGoogleAuthConfig()
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(models.MessageHTTP{Message: err.Error()})
	}
	tokenAuth2, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.MessageHTTP{Message: err.Error()})
	}
	calendarRepository := repositories.NewCalendarRepository(h.db)
	err = calendarRepository.InsertGoogleCalendarToken(tokenAuth2, userId)
	if err != nil {
		log.Print(err)
		return fiber.ErrBadGateway
	}
	return c.SendFile("./public/google_token_registred.html")
}

func (h Handler) GetEventList(c *fiber.Ctx) error {
	srv, httpModelError := services.
		NewCalendarServiceFactor(h.db).
		GetCalendarService(
			c.Locals(constants.Token).(string),
			c.Locals(constants.UserId),
		)
	if httpModelError != nil {
		return c.Status(httpModelError.HttpCode).JSON(httpModelError)
	}
	initialTime, err := time.Parse(time.RFC3339, c.Query(constants.InitialDate))
	if err != nil {
		log.Printf("Deu ruim %v", err)
		return models.
			MessageHTTPFromFiberError(fiber.ErrBadRequest).
			FiberContext(c)
	}
	events, err := srv.Events.
		List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(initialTime.Format(time.RFC3339)).
		TimeMax(initialTime.Add(time.Hour * 24).Format(time.RFC3339)).
		OrderBy("startTime").
		Do()
	if err != nil {
		log.Printf("Unable to retrieve next ten of the user's events: %v", err)
		return fiber.ErrInternalServerError
	}
	return c.Status(200).JSON(events.Items)
}
