package repositories

import (
	"database/sql"
	"errors"
	"golang.org/x/oauth2"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return AuthRepository{db}
}

func (a *AuthRepository) GetToken(userId int) (*oauth2.Token, error) {
	var token oauth2.Token
	var expiry string
	err := a.db.QueryRow(
		"select access_token, token_type, refresh_token, expiry from google_calendar_token where user_id = $1",
		userId,
	).Scan(&token.AccessToken, &token.TokenType, &token.RefreshToken, &expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, TokenNotFounded
		}
		return nil, err
	}
	token.Expiry, err = time.Parse("2006-01-02 15:04:05.999999999Z07:00", expiry)
	if err != nil {
		return nil, err
	}
	return &token, err
}

func (a *AuthRepository) DeleteCalendarTokenByUserId(userId any) error {
	_, err := a.db.Exec("delete from google_calendar_token where user_id = $1", userId)
	return err
}

func (a *AuthRepository) InsertTokenBlackList(userId int, token string, expiry time.Time) error {
	_, err := a.db.Exec(
		"insert into token_black_list values ($1, $2, $3)",
		userId,
		token,
		expiry,
	)
	return err
}

func (a *AuthRepository) UpdateToken(id int, token *oauth2.Token) error {
	_, err := a.db.Exec(
		"update google_calendar_token set access_token = $2, expiry = $3 where user_id = $1",
		id,
		token.AccessToken,
		token.Expiry,
	)
	return err
}

func (a *AuthRepository) IsValidToken(tokenString string) (bool, error) {
	count := 0
	err := a.db.
		QueryRow("select count(*) from token_black_list where token = $1", tokenString).
		Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
