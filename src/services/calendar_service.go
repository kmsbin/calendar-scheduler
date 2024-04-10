package services

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"context"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

type CalendarServiceFactory struct {
	db      *sql.DB
	baseUrl string
}

func NewCalendarServiceFactor(db *sql.DB, baseUrl string) CalendarServiceFactory {
	return CalendarServiceFactory{db, baseUrl}
}

func (calendarServiceFactory CalendarServiceFactory) GetCalendarServiceByContext(c *fiber.Ctx) (*calendar.Service, *models.MessageHTTP) {
	token := c.Locals(constants.Token).(string)
	userId := c.Locals(constants.UserId)
	return calendarServiceFactory.GetCalendarService(token, userId)
}

func (calendarServiceFactory CalendarServiceFactory) GetCalendarServiceByUserId(userId any) (*calendar.Service, *models.MessageHTTP) {
	return calendarServiceFactory.GetCalendarService("", userId)
}

func (calendarServiceFactory CalendarServiceFactory) GetCalendarService(token string, userId any) (*calendar.Service, *models.MessageHTTP) {
	if userId == nil {
		return nil, models.MessageHTTPFromFiberError(fiber.ErrUnauthorized)
	}
	userRepository := repositories.NewUserRepository(calendarServiceFactory.db)
	user, err := userRepository.GetUserById(userId)
	if err != nil {
		return nil, &models.MessageHTTP{Message: "User not founded!", HttpCode: fiber.StatusUnauthorized}
	}
	config := repositories.NewGoogleCalendarRepository(token, "").GetGoogleAuthConfig()
	client, err := calendarServiceFactory.getClient(token, user.Id, config)

	if err != nil {
		if tokenNotFoundedErr, ok := err.(CalendarTokenNotFounded); ok {
			return nil, &models.MessageHTTP{
				Message:  tokenNotFoundedErr.AuthUrl,
				HttpCode: fiber.StatusPreconditionRequired,
			}
		}
		return nil, models.MessageHTTPFromFiberError(fiber.ErrInternalServerError)
	}
	service, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, &models.MessageHTTP{Message: err.Error(), HttpCode: fiber.StatusInternalServerError}
	}
	return service, nil
}

func (calendarServiceFactory CalendarServiceFactory) getClient(token string, userId int, config *oauth2.Config) (*http.Client, error) {
	tok, err := calendarServiceFactory.tokenFromDb(userId, config)
	if err != nil {
		return nil, CalendarTokenNotFounded{config.AuthCodeURL(
			"state-token",
			oauth2.AccessTypeOffline,
			oauth2.SetAuthURLParam("state", token),
		)}
	}

	return config.Client(context.Background(), tok), nil
}

func (calendarServiceFactory CalendarServiceFactory) deleteGoogleAuthToken(userId any, authRepository repositories.AuthRepository) error {
	err := authRepository.DeleteCalendarTokenByUserId(userId)
	if err != nil {
		log.Fatalln(err)
	}
	return errors.New("google credential is expired")
}

func (calendarServiceFactory CalendarServiceFactory) tokenFromDb(userId int, config *oauth2.Config) (*oauth2.Token, error) {
	authRepository := repositories.NewAuthRepository(calendarServiceFactory.db)
	token, err := authRepository.GetToken(userId)
	if err != nil {
		log.Printf("erro %v", err.Error())
		return nil, err
	}
	tokenReuse := config.TokenSource(context.TODO(), token)
	newToken, err := tokenReuse.Token()
	if err != nil {
		if err := calendarServiceFactory.deleteGoogleAuthToken(userId, authRepository); err != nil {
			return nil, err
		}
		return nil, err
	}
	if newToken.AccessToken != token.AccessToken {
		if err = authRepository.UpdateToken(userId, newToken); err != nil {
			return nil, calendarServiceFactory.deleteGoogleAuthToken(userId, authRepository)
		}
		token = newToken
	}
	if &token.AccessToken == nil {
		return nil, errors.New("token not founded")
	}
	return token, nil
}

type CalendarTokenNotFounded struct {
	AuthUrl string
}

func (ctf CalendarTokenNotFounded) Error() string {
	return ctf.AuthUrl
}
