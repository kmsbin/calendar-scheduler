package handlers

import (
	constants "calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func (h Handler) GetEventsFree(c *fiber.Ctx) error {
	meetingsRange, err := h.getmeetingsRangeFromCode(c.Query(constants.Code))
	if err != nil {
		return err.FiberContext(c)
	}
	return c.JSON(meetingsRange)
}

func (h Handler) getmeetingsRangeFromCode(code string) (*models.MeetingsRange, *models.MessageHTTP) {
	if len(code) == 0 {
		return nil, &models.MessageHTTP{
			Message:  "missing query parameter",
			HttpCode: fiber.StatusUnprocessableEntity,
		}
	}
	meetingsRepository := repositories.NewMeetingsRepository(h.db)
	meetingsRange, err := meetingsRepository.GetmeetingsRangeByCode(code)
	if err != nil {
		if errors.Is(err, repositories.MeetingsRangeNotFounded) {
			return nil, &models.MessageHTTP{
				Message:  err.Error(),
				HttpCode: fiber.StatusNotFound,
			}
		}
		return nil, models.MessageHTTPFromFiberError(fiber.ErrInternalServerError)

	}
	return meetingsRange, nil
}
