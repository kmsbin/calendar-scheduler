package repositories

import (
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"os"
)

type GoogleCalendarRepository struct {
	token   string
	baseUrl string
}

func NewGoogleCalendarRepository(token, baseUrl string) GoogleCalendarRepository {
	return GoogleCalendarRepository{
		token:   token,
		baseUrl: baseUrl,
	}
}

func (g GoogleCalendarRepository) GetGoogleAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_AUTH_CLIENT_ID"), //"536852565523-rnm7n3d4bj00uu5tmb513bls2ucpepev.apps.googleusercontent.com",
		ClientSecret: os.Getenv("GOOGLE_AUTH_CLIENT_SECRET"),
		RedirectURL:  fmt.Sprintf("%s/v1/set-token-google", g.baseUrl),
		Scopes:       []string{calendar.CalendarEventsScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
}

func (g GoogleCalendarRepository) GetGoogleAuthUrl(token string) string {
	config := g.GetGoogleAuthConfig()
	authUrl := config.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("state", token),
	)
	return authUrl
}
