package handlers

import (
	"calendar_scheduler/src"
	"calendar_scheduler/src/database"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/calendar/v3"

	"log"
	"time"
)

func SetTokenCalendar(c *fiber.Ctx) error {
	token := c.Query("state")
	code := c.Query("code")
	if len(token) == 0 || len(code) == 0 {
		log.Printf("state %s", token)
		return fiber.ErrUnauthorized
	}

	err := ValidateToken(token, c)
	if err != nil {
		log.Printf("Erro doido %v", err)
		return err
	}
	log.Printf("state redirect ")
	userId := c.Locals("user_id")

	if userId == nil {
		return fiber.ErrUnauthorized
	}
	_, err = repositories.NewUserRepository().GetDataFromToken(userId.(float64))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.MessageHTTP{Message: "User not founded!"})
	}
	config, err := repositories.GetGoogleAuthConfig()
	if err != nil {
		msg := models.MessageHTTP{Message: err.Error()}
		return c.Status(fiber.StatusInternalServerError).JSON(msg)
	}
	tokenAuth2, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.MessageHTTP{Message: err.Error()})
	}

	db, _ := database.OpenConnection()
	_, err = db.Exec("insert into google_calendar_token(user_id, access_token, refresh_token, expiry, token_type) values ($1, $2, $3, $4, $5)",
		userId, tokenAuth2.AccessToken, tokenAuth2.RefreshToken, tokenAuth2.Expiry, tokenAuth2.TokenType)
	if err != nil {
		log.Print(err)
		return fiber.ErrBadGateway
	}
	return c.SendFile("./public/google_token_registred.html")
}

func GetEventList(c *fiber.Ctx) error {
	srv, httpModelError := GetCalendarService(c)
	if httpModelError != nil {
		return c.Status(httpModelError.HttpCode).JSON(httpModelError)
	}
	initialTime, err := time.Parse(time.RFC3339, c.Query("initial_date"))
	if err != nil {
		log.Printf("Deu ruim %v", err)
		return models.MessageHTTPFromFiberError(fiber.ErrBadRequest)
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

func CreateEvent(c *fiber.Ctx) error {
	calendarEvent := models.CalendarEvent{}
	if err := c.BodyParser(&calendarEvent); err != nil {
		log.Printf("Unable to parse event. %v\n", err)
		return models.MessageHTTPFromFiberError(fiber.ErrUnprocessableEntity)
	}
	srv, httpModelError := GetCalendarService(c)
	if httpModelError != nil {
		return c.Status(httpModelError.HttpCode).JSON(httpModelError)
	}

	event := &calendar.Event{
		Summary:     calendarEvent.Summary,
		Description: calendarEvent.Description,
		Start: &calendar.EventDateTime{
			DateTime: calendarEvent.Start.Format(time.RFC3339),
			TimeZone: src.Locale,
		},
		End: &calendar.EventDateTime{
			DateTime: calendarEvent.End.Format(time.RFC3339),
			TimeZone: src.Locale,
		},
		Attendees: []*calendar.EventAttendee{
			{Email: calendarEvent.Email},
		},
	}
	event, err := srv.Events.Insert(src.CalendarId, event).Do()
	if err != nil {
		log.Printf("Unable to create event. %v\n", err)
		return fiber.ErrInternalServerError
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
	return c.Status(200).JSON(event)
}
