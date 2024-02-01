package handlers

import (
	"calendar_scheduler/src/repositories"
	"context"
	"errors"
	"golang.org/x/oauth2"
	"log"
	"net/http"
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
	authRepository := repositories.NewAuthRepository()
	token, err := authRepository.GetToken(userId)
	if !token.Valid() {
		err = authRepository.DeleteTokenByUserId(userId)
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
	return token, nil
}

type calendarTokenNotFounded struct {
	AuthUrl string
}

func (ctf calendarTokenNotFounded) Error() string {
	return ctf.AuthUrl
}
