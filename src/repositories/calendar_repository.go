package repositories

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"log"
)

type CalendarRepository struct {
	db *sql.DB
}

func NewCalendarRepository(db *sql.DB) CalendarRepository {
	return CalendarRepository{db}
}

func (c *CalendarRepository) InsertGoogleCalendarToken(tokenAuth2 *oauth2.Token, userId any) error {
	_, err := c.db.Exec("insert into google_calendar_token(user_id, access_token, refresh_token, expiry, token_type) values ($1, $2, $3, $4, $5)",
		userId,
		tokenAuth2.AccessToken,
		tokenAuth2.RefreshToken,
		tokenAuth2.Expiry,
		tokenAuth2.TokenType,
	)
	if err != nil {
		log.Print(err)
		return fiber.ErrBadGateway
	}
	return nil
}
