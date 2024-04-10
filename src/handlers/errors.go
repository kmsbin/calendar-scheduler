package handlers

import (
	"calendar_scheduler/src/models"
	"github.com/gofiber/fiber/v2"
)

func InternalServerError(c *fiber.Ctx, data ...any) error {
	return errorBuilder(c, fiber.ErrInternalServerError, data...)
}

func BadGatewayError(c *fiber.Ctx, data ...any) error {
	return errorBuilder(c, fiber.ErrBadGateway, data...)
}

func NotFoundError(c *fiber.Ctx, data ...any) error {
	return errorBuilder(c, fiber.ErrNotFound, data...)
}

func UnauthorizedError(c *fiber.Ctx, data ...any) error {
	return errorBuilder(c, fiber.ErrUnauthorized, data...)
}

func NotAcceptableError(c *fiber.Ctx, data ...any) error {
	return errorBuilder(c, fiber.ErrNotAcceptable, data...)
}

func UnprocessableEntity(c *fiber.Ctx, data ...any) error {
	return errorBuilder(c, fiber.ErrUnprocessableEntity, data...)
}

func BadRequestError(c *fiber.Ctx, data ...any) error {
	return errorBuilder(c, fiber.ErrBadRequest, data...)
}

func ConflictError(c *fiber.Ctx, data ...any) error {
	return errorBuilder(c, fiber.ErrConflict, data...)
}

func GoneError(c *fiber.Ctx, data ...any) error {
	return errorBuilder(c, fiber.ErrGone, data...)
}

func ResponseOK(c *fiber.Ctx, data ...any) error {
	c.Status(200)
	switch len(data) {
	case 0:
		return c.JSON(models.MessageHTTPFromMessage("Successful!"))
	case 1:
		return c.JSON(data[0])
	default:
		return c.JSON(data)
	}
}

func errorBuilder(c *fiber.Ctx, err *fiber.Error, data ...any) error {
	c.Status(err.Code)
	switch len(data) {
	case 0:
		return c.JSON(models.MessageHTTPFromMessage(err.Message))
	case 1:
		return c.JSON(data[0])
	default:
		return c.JSON(data)
	}
}
