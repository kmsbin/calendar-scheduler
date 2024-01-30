package handlers

import (
	"calendar_scheduler/src/database"
	"context"
	"errors"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

func getClient(token string, userId int, config *oauth2.Config) (*http.Client, error) {
	tok, err := tokenFromDb(userId)
	if err == nil {
		return config.Client(context.Background(), tok), nil
	}
	log.Printf("getClient err %v", err)
	return nil, calendarTokenNotFounded{config.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("state", token),
	)}
}

func tokenFromDb(userId int) (*oauth2.Token, error) {
	db, _ := database.OpenConnection()
	var token oauth2.Token
	var expiry string
	err := db.QueryRow(
		"select access_token, token_type, refresh_token, expiry from google_calendar_token where user_id = $1",
		userId,
	).Scan(&token.AccessToken, &token.TokenType, &token.RefreshToken, &expiry)
	token.Expiry, _ = time.Parse("2006-01-02 15:04:05.999999999Z07:00", expiry)
	if !token.Valid() {
		_, err = db.Exec("delete from google_calendar_token where user_id = $1", userId)
		if err != nil {
			panic(err)
		}
		return nil, errors.New("google credential is expired")
	}

	if err != nil {
		log.Printf("erro %v", err.Error())
		return nil, err
	}
	if &token.AccessToken == nil {
		return nil, errors.New("token not founded")
	}
	return &token, nil
}

type calendarTokenNotFounded struct {
	AuthUrl string
}

func (ctf calendarTokenNotFounded) Error() string {
	return ctf.AuthUrl
}
