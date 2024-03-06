package repositories

import (
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"log"
	"os"
)

func GetGoogleAuthConfig() (*oauth2.Config, error) {
	credentials, err := os.ReadFile("credentials/credentials.json")
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		return nil, errors.New("cannot read calendar token")
	}
	config, err := google.ConfigFromJSON(credentials, calendar.CalendarEventsScope)
	if err != nil {
		log.Printf("Unable create google config token: %v", err)
		return nil, errors.New("unable create google config token")
	}
	return config, nil
}

func GetGoogleAuthUrl(token string) (string, error) {
	config, err := GetGoogleAuthConfig()
	if err != nil {
		return "", nil
	}
	authUrl := config.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("state", token),
	)
	return authUrl, nil
}
