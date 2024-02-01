package main

import (
	"calendar_scheduler/src"
	"calendar_scheduler/src/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"time"
)

func main() {
	time.Local, _ = time.LoadLocation(src.Locale)
	app := fiber.New()
	app.Use(logger.New())

	app.Post("/users", handlers.CreateUserHadler)
	app.Get("/sign-in", handlers.SignInUser)
	app.Get("/set-token-google", handlers.SetTokenCalendar)
	app.Use("/app", handlers.ValidateTokenMiddleware)
	app.Get("/app/user", handlers.GetUser)
	app.Get("/app/events", handlers.GetEventList)
	app.Get("/app/events-free", handlers.GetEmptyScheduledTime)
	app.Post("/app/event", handlers.CreateEvent)
	app.Post("/app/meeting-range", handlers.CreateMeetingRange)
	app.Get("/app/meeting-range", handlers.GetMeetingRange)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
