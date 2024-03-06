package main

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/database"
	"calendar_scheduler/src/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"time"
)

func main() {
	time.Local, _ = time.LoadLocation(constants.Locale)
	dbConn, err := database.OpenConnection()
	if err != nil {
		log.Fatalf("Error opening db connection %v", err)
	}
	defer database.CloseConnection(dbConn)
	handler := handlers.NewHandler(dbConn)
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Post("/sign-up", handler.CreateUserHadler)
	app.Get("/sign-in", handler.SignInUser)
	app.Get("/set-token-google", handler.SetTokenCalendar)
	app.Get("/get-events-free", handler.GetEventsFree)
	app.Get("/get-events-code", handler.GetEventsByCode)
	app.Post("/event", handler.CreateGoogleCalendarEvent)
	app.Use("/app", handler.ValidateTokenMiddleware)
	app.Get("/app/user", handler.GetUser)
	app.Delete("/app/delete-user", handler.DeleteUser)
	app.Get("/app/events", handler.GetEventList)
	app.Post("/app/meetings-range", handler.CreatemeetingsRange)
	app.Get("/app/meetings-range", handler.GetmeetingsRange)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
	err = app.Listen(":3000")
	if err != nil {
		log.Fatalln(err)
	}
}
