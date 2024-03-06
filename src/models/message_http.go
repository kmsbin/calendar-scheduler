package models

import (
	"github.com/gofiber/fiber/v2"
)

type MessageHTTP struct {
	Code     int             `json:"code,omitempty"`
	HttpCode int             `json:"-"`
	Message  string          `json:"message"`
	Extra    *map[string]any `json:"-"`
}

func MessageHTTPFromFiberError(err *fiber.Error) *MessageHTTP {
	return &MessageHTTP{
		HttpCode: err.Code,
		Message:  err.Message,
	}
}

func MessageHTTPFromMessage(message string) *MessageHTTP {
	return &MessageHTTP{Message: message}
}

func (m *MessageHTTP) FiberContext(c *fiber.Ctx) error {
	return c.Status(m.HttpCode).JSON(m)
}
