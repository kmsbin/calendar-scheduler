package main

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/database"
	"calendar_scheduler/src/handlers"
	"calendar_scheduler/src/helpers"
	"database/sql"
	_ "github.com/joho/godotenv/autoload"

	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
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
	if helpers.IsRunningInLambda() {
		lambda.Start(fiberadapter.New(app).ProxyWithContext)
	} else {
		log.Fatal(app.Listen(":3000"))
	}
}

func setHandlers(app *fiber.App, dbConn *sql.DB) {
	handler := handlers.NewHandler(dbConn)
	//app.Use(filesystem.New(filesystem.Config{
	//	Root: http.FS(static),
	//}))
	app.Get("/well-know", func(c *fiber.Ctx) error {
		return c.Status(200).JSON("tudo certo chapa")
	})
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
	app.Use("/api", handler.ValidateTokenMiddleware)
	// user info
	app.Get("/api/user", handler.GetUser)
	app.Delete("/api/delete-user", handler.DeleteUser)
	app.Get("/api/events", handler.GetEventList)
	app.Post("/api/meeting-range", handler.CreateMeetingsRange)
	app.Get("/api/meeting-range", handler.GetMeetingsRange)
	//app.Static("/web", "./public/web/index.html")
	//app.Static("/app/", "./public/app/")
	app.Use(func(c *fiber.Ctx) error {
		log.Println("base url", c.BaseURL())
		//if path := c.Path(); strings.HasPrefix(path, "/app") || path == "/" {
		//	return c.SendFile("./public/app/")
		//}
		return c.SendStatus(404)
	})
}

func setFiberConfigs(app *fiber.App) {
	app.Use(logger.New())

	allowedMethods := []string{
		fiber.MethodHead,
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodDelete,
		fiber.MethodOptions,
	}

	corsConfig := cors.Config{
		AllowMethods: strings.Join(allowedMethods, ","),
		AllowOrigins: "*",
	}
	app.Use(cors.New(corsConfig))
	app.Use(helmet.New())
}
