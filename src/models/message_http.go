package models

import "github.com/gofiber/fiber/v2"

type MessageHTTP struct {
	Code     int    `json:"code,omitempty"`
	HttpCode int    `json:"-"`
	Message  string `json:"message"`
}

func (m *MessageHTTP) Error() string {
	return m.Message
}

func MessageHTTPFromFiberError(err *fiber.Error) *MessageHTTP {
	return &MessageHTTP{
		HttpCode: err.Code,
		Code:     err.Code,
		Message:  err.Message,
	}
}
