package handlers

import (
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
)

func CreateMeetingRange(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id")
	if userId == nil {
		log.Println("User id is nil.")
		return models.MessageHTTPFromFiberError(fiber.ErrUnauthorized)
	}
	meetingBody := models.MeetingRange{
		UserId: int(userId.(float64)),
	}
	if err := ctx.BodyParser(&meetingBody); err != nil {
		return models.MessageHTTPFromFiberError(fiber.ErrUnprocessableEntity)
	}
	if err := validateMeetingRange(meetingBody); err != nil {
		return ctx.
			Status(fiber.StatusUnprocessableEntity).
			JSON(models.MessageHTTP{
				HttpCode: fiber.StatusUnprocessableEntity,
				Message:  err.Error(),
			})
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
	meetingRange, err := meetingRepository.GetLastMeetingRange(userId)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return ctx.Status(fiber.StatusOK).JSON(meetingRange)
}

func validateMeetingRange(meetingRange models.MeetingRange) error {
	start, end, err := meetingRange.ConvertToTime()
	if err != nil {
		log.Printf("dates %v, %v", meetingRange.Start, meetingRange.End)
		log.Printf("errors: %v", err)
		return errors.New("error parsing dates")
	}
	if start.After(*end) {
		return errors.New("the start date cannot be after end date")
	}
	return nil
}
