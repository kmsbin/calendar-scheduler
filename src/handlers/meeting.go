package handlers

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
)

func (h Handler) CreateMeetingsRange(ctx *fiber.Ctx) error {
	userId, token := ctx.Locals(constants.UserId), ctx.Locals(constants.Token).(string)
	if userId == nil {
		return models.
			MessageHTTPFromFiberError(fiber.ErrUnauthorized).
			FiberContext(ctx)
	}
	authUrl, err := repositories.GetGoogleAuthUrl(token)
	if err != nil {
		return models.MessageHTTPFromFiberError(fiber.ErrInternalServerError).FiberContext(ctx)
	}

	meetingsBody := models.MeetingsRange{
		UserId: int(userId.(float64)),
	}
	if err := ctx.BodyParser(&meetingsBody); err != nil {
		return models.
			MessageHTTPFromFiberError(fiber.ErrUnprocessableEntity).
			FiberContext(ctx)
	}
	if err := validatemeetingsRange(meetingsBody); err != nil {
		return ctx.
			Status(fiber.StatusUnprocessableEntity).
			JSON(models.MessageHTTPFromMessage(err.Error()))
	}
	meetingsRepository := repositories.NewMeetingsRepository(h.db)
	meetingsBody.Code = uuid.New().String()
	err = meetingsRepository.InsertMeetingsRange(meetingsBody)
	if err != nil {
		log.Print(err)
		return models.
			MessageHTTPFromFiberError(fiber.ErrBadGateway).FiberContext(ctx)
	}
	return ctx.
		Status(fiber.StatusOK).
		JSON(createmeetingsRangeResponse{authUrl})
}

func (h Handler) GetMeetingsRange(c *fiber.Ctx) error {
	userId := c.Locals(constants.UserId)
	if userId == nil {
		log.Println("User id is nil.")
		return models.
			MessageHTTPFromFiberError(fiber.ErrUnauthorized).
			FiberContext(c)
	}
	meetingsRepository := repositories.NewMeetingsRepository(h.db)
	meetingsRange, err := meetingsRepository.GetLastmeetingsRange(userId)
	if err != nil {
		log.Println(err)
		if errors.Is(err, repositories.MeetingsRangeNotFounded) {
			return models.
				MessageHTTPFromFiberError(fiber.ErrNotFound).
				FiberContext(c)
		}
		return models.
			MessageHTTPFromFiberError(fiber.ErrInternalServerError).
			FiberContext(c)
	}
	return c.
		Status(fiber.StatusOK).
		JSON(meetingsRange)
}

func validatemeetingsRange(meetingsRange models.MeetingsRange) error {
	start, end, err := meetingsRange.ConvertToTime()
	if err != nil {
		log.Printf("dates %v, %v", meetingsRange.Start, meetingsRange.End)
		log.Printf("errors: %v", err)
		return errors.New("error parsing dates")
	}
	if start.After(*end) {
		return errors.New("the start date cannot be after end date")
	}
	return nil
}

type createmeetingsRangeResponse struct {
	Url string `json:"url"`
}
