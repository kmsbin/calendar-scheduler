package handlers

import (
	constants "calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func (h Handler) GetEventsFree(c *fiber.Ctx) error {
	meetingsRange, err := h.getmeetingsRangeFromCode(c)
	if err != nil {
		return err
	}
	return c.JSON(meetingsRange)
}

func (h Handler) getmeetingsRangeFromCode(c *fiber.Ctx) (*models.meetingsRange, error) {
	code := c.Query(constants.Code)
	if len(code) == 0 {
		return nil, c.
			Status(fiber.StatusUnprocessableEntity).
			JSON(models.MessageHTTP{Message: "missing query parameter"})
	}
	meetingsRepository := repositories.NewmeetingsRepository(h.db)
	meetingsRange, err := meetingsRepository.GetmeetingsRangeByCode(code)
	if err != nil {
		if errors.Is(err, repositories.meetingsRangeNotFounded) {
			return nil, c.
				Status(fiber.StatusNotFound).
				JSON(models.MessageHTTP{Message: err.Error()})
		}
		return nil, c.
			Status(fiber.StatusInternalServerError).
			JSON(models.MessageHTTP{Message: fiber.ErrInternalServerError.Error()})

	}
	return meetingsRange, nil
}
