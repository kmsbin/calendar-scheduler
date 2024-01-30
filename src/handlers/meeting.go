package handlers

import (
	"calendar_scheduler/src/database"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"github.com/gofiber/fiber/v2"
	"log"
	time "time"
)

func CreateMeeting(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id")
	if userId == nil {
		return models.MessageHTTPFromFiberError(fiber.ErrUnauthorized)
	}
	meetingBody := models.Meeting{}
	if err := ctx.BodyParser(&meetingBody); err != nil {
		log.Printf("Unable to parse event. %v\n", err)
		return models.MessageHTTPFromFiberError(fiber.ErrUnprocessableEntity)
	}
	if err := validateMeeting(&meetingBody); err != nil {
		return err
	}
	db, _ := database.OpenConnection()
	_, err := db.Exec(
		"insert into meeting(user_id, summary, start_date, end_date, duration) values ($1, $2, $3, $4, $5)",
		userId,
		meetingBody.Summary,
		meetingBody.Start.ToString(),
		meetingBody.End.ToString(),
		meetingBody.Duration,
	)
	if err != nil {
		log.Print(err)
		return models.MessageHTTPFromFiberError(fiber.ErrBadGateway)
	}
	return ctx.JSON(models.MessageHTTP{
		HttpCode: fiber.StatusCreated,
		Message:  "Meeting created",
	})
}

func CreateMeetingRange(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id")
	if userId == nil {
		log.Println("User id is nil.")
		return models.MessageHTTPFromFiberError(fiber.ErrUnauthorized)
	}
	meetingBody := models.MeetingRange{
		UserId: userId.(int),
	}
	if err := ctx.BodyParser(&meetingBody); err != nil {
		log.Printf("Unable to parse event. %v\n", err)
		return models.MessageHTTPFromFiberError(fiber.ErrUnprocessableEntity)
	}
	meetingRepository := repositories.NewMeetingRepository()
	err := meetingRepository.InsertMeetingRange(meetingBody)
	if err != nil {
		log.Print(err)
		return models.MessageHTTPFromFiberError(fiber.ErrBadGateway)
	}
	return ctx.JSON(models.MessageHTTP{
		HttpCode: fiber.StatusCreated,
		Message:  "Meeting created",
	})
}

func GetMeetingRange(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id")
	if userId == nil {
		log.Println("User id is nil.")
		return models.MessageHTTPFromFiberError(fiber.ErrUnauthorized)
	}
	meetingRepository := repositories.NewMeetingRepository()
	meetingRange, err := meetingRepository.GetLastMeetingRange(userId.(int))
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return ctx.Status(fiber.StatusOK).JSON(meetingRange)
}

func validateMeeting(meeting *models.Meeting) error {
	if meeting.Start.After(meeting.End.Time) {
		return &models.MessageHTTP{
			HttpCode: fiber.StatusUnprocessableEntity,
			Message:  "The start date cannot be after the end date",
		}
	}
	if meeting.Start.Before(time.Now()) || meeting.End.Before(time.Now()) {
		return &models.MessageHTTP{
			HttpCode: fiber.StatusUnprocessableEntity,
			Message:  "The date cannot be before the current date",
		}
	}
	return nil
}
