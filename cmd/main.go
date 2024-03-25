package main

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/database"
	"calendar_scheduler/src/handlers"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"strings"
	"time"
)

func main() {
	time.Local, _ = time.LoadLocation(constants.Locale)
	dbConn, err := database.OpenConnection()
	if err != nil {
		log.Fatalf("Error opening db connection %v", err)
	}
	defer database.CloseConnection(dbConn)

	app := fiber.New()
	app.Use(logger.New())

	setFiberConfigs(app)
	setHandlers(app, dbConn)
	err = app.Listen(":3000")
	if err != nil {
		log.Fatalln(err)
	}
}

func setHandlers(app *fiber.App, dbConn *sql.DB) {
	handler := handlers.NewHandler(dbConn)
	// Auth
	app.Get("/sign-in", handler.SignInUser)
	app.Post("/sign-up", handler.SignUpUser)
	app.Delete("/sign-out", handler.SignOutUser)
	app.Get("/send-password-recover", handler.SendPasswordRecover)
	app.Post("/receive-password-recover", handler.ReceivePasswordRecover)
	// Google auth sign in
	app.Get("/set-token-google", handler.SetTokenGoogleCalendar)
	// calendar info
	app.Get("/get-events-code", handler.GetEventsByCode)
	app.Post("/event", handler.CreateGoogleCalendarEvent)
	app.Use("/app", handler.ValidateTokenMiddleware)
	app.Get("/app/user", handler.GetUser)
	app.Delete("/app/delete-user", handler.DeleteUser)
	app.Get("/app/events", handler.GetEventList)
	app.Post("/app/meeting-range", handler.CreateMeetingsRange)
	app.Get("/app/meeting-range", handler.GetMeetingsRange)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
}

func setFiberConfigs(app *fiber.App) {
	app.Use(logger.New())

	allowedMethods := []string{
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodDelete,
	}

	corsConfig := cors.Config{
		AllowMethods: strings.Join(allowedMethods, ","),
	}
	app.Use(cors.New(corsConfig))
	app.Use(helmet.New())
}
