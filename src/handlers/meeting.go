package handlers

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"time"
)

func (h Handler) CreateMeetingsRange(c *fiber.Ctx) error {
	userId, token := c.Locals(constants.UserId), c.Locals(constants.Token).(string)
	if userId == nil {
		return UnauthorizedError(c)
	}
	authUrl := repositories.NewGoogleCalendarRepository(token, c.BaseURL()).GetGoogleAuthUrl(token)

	meetingsBody := models.MeetingsRange{
		UserId: int(userId.(float64)),
	}
	if err := c.BodyParser(&meetingsBody); err != nil {
		return UnprocessableEntity(c)
	}
	if err := validateMeetingsRange(meetingsBody); err != nil {
		return UnprocessableEntity(c, models.MessageHTTPFromMessage(err.Error()))
	}
	meetingsRepository := repositories.NewMeetingsRepository(h.db)
	meetingsBody.Code = uuid.New().String()
	err := meetingsRepository.InsertMeetingsRange(meetingsBody)
	if err != nil {
		log.Print(err)
		return BadGatewayError(c)
	}

	return ResponseOK(c, createmeetingsRangeResponse{authUrl})
}

func (h Handler) GetMeetingsRange(c *fiber.Ctx) error {
	userId := c.Locals(constants.UserId)
	token := c.Locals(constants.Token).(string)
	if userId == nil {
		log.Println("User id is nil.")
		return UnauthorizedError(c)
	}
	meetingsRepository := repositories.NewMeetingsRepository(h.db)
	meetingsRange, err := meetingsRepository.GetLastMeetingsRange(userId)
	if err != nil {
		log.Println(err)
		if errors.Is(err, repositories.MeetingsRangeNotFounded) {
			return NotFoundError(c)
		}
		return InternalServerError(c)
	}
	authRepository := repositories.NewAuthRepository(h.db)
	googleToken, err := authRepository.GetToken(int(userId.(float64)))
	if err != nil && !errors.Is(err, repositories.TokenNotFounded) {
		log.Println(err)
		return InternalServerError(c)
	}
	if googleToken == nil {
		meetingsRange.AuthUrl = repositories.
			NewGoogleCalendarRepository(token, c.BaseURL()).
			GetGoogleAuthUrl(token)
	}
	return ResponseOK(c, meetingsRange)
}

func validateMeetingsRange(meetingsRange models.MeetingsRange) error {
	start, end, err := meetingsRange.ConvertToTime()
	if err != nil {
		log.Printf("dates %v, %v", meetingsRange.Start, meetingsRange.End)
		log.Printf("errors: %v", err)
		return errors.New("error parsing dates")
	}
	if _, err := time.ParseDuration(meetingsRange.Duration); err != nil {
		return errors.New("error parsing duration")
	}
	if start.After(*end) {
		return errors.New("the start date cannot be after end date")
	}
	return nil
}

type createmeetingsRangeResponse struct {
	Url string `json:"url"`
}
