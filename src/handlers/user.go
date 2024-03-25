package handlers

import (
	"calendar_scheduler/src/auth"
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func (h Handler) SignUpUser(c *fiber.Ctx) error {
	c.Set(headers.ContentType, fiber.MIMEApplicationJSON)
	var user models.User
	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		log.Println(err)
		return models.MessageHTTPFromFiberError(fiber.ErrNotAcceptable).FiberContext(c)
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return models.MessageHTTPFromFiberError(fiber.ErrInternalServerError).FiberContext(c)
	}
	userRepository := repositories.NewUserRepository(h.db)
	err = userRepository.InsertUser(&user, password)
	if err != nil {
		log.Println(err)
		return models.MessageHTTPFromFiberError(fiber.ErrInternalServerError).FiberContext(c)
	}
	return c.
		Status(fiber.StatusOK).
		JSON(models.MessageHTTPFromMessage("Successful!!"))
}

func (h Handler) GetUser(c *fiber.Ctx) error {
	userId := c.Locals(constants.UserId)
	if userId == nil {
		return models.MessageHTTPFromFiberError(fiber.ErrUnauthorized).FiberContext(c)
	}
	userRepository := repositories.NewUserRepository(h.db)
	user, err := userRepository.GetUserById(userId)
	if err != nil {
		return models.
			MessageHTTPFromFiberError(fiber.ErrNotFound).
			FiberContext(c)
	}
	return c.
		Status(fiber.StatusOK).
		JSON(user)
}
func (h Handler) SignInUser(c *fiber.Ctx) error {
	email, password := c.Query("email"), c.Query("password")
	if email == "" || password == "" {
		return fiber.ErrNotAcceptable
	}
	userRepository := repositories.NewUserRepository(h.db)
	user, userPassword, err := userRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.UserNotFounded) {
			return c.
				Status(fiber.StatusNotFound).
				JSON(models.MessageHTTPFromFiberError(fiber.ErrNotFound))
		}
		return models.
			MessageHTTPFromFiberError(fiber.ErrInternalServerError).
			FiberContext(c)
	}
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	if err != nil {
		return models.
			MessageHTTPFromFiberError(fiber.ErrUnauthorized).
			FiberContext(c)
	}
	token, err := auth.CreateToken(*user)
	if err != nil {
		return models.
			MessageHTTPFromFiberError(fiber.ErrUnauthorized).
			FiberContext(c)
	}
	return c.
		Status(fiber.StatusOK).
		JSON(map[string]string{constants.Token: token})
}

func (h Handler) ValidateTokenMiddleware(c *fiber.Ctx) error {
	tokenString := c.Query(constants.Token)
	if httpModel := ValidateToken(tokenString, c); httpModel != nil {
		return httpModel.FiberContext(c)
	}

	authRepository := repositories.NewAuthRepository(h.db)
	isValid, err := authRepository.IsValidToken(tokenString)
	if err != nil {
		log.Printf("error validating token %v \n", err)
		return models.MessageHTTPFromFiberError(fiber.ErrInternalServerError).FiberContext(c)
	}
	if !isValid {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(models.MessageHTTPFromMessage("Invalid token"))
	}
	if err := c.Next(); err != nil {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(models.MessageHTTPFromMessage("Invalid token"))
	}
	return nil
}

func validateToken(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(auth.SigningKey), nil
}

func GetTokenExpirationData(tokenString string, c *fiber.Ctx) (*TokenData, *models.MessageHTTP) {
	token, err := jwt.Parse(tokenString, validateToken)

	if httpError := normalizeJwtError(token, err); httpError != nil {
		return nil, httpError
	}
	expiration, err := time.Parse(time.RFC3339, token.Claims.(jwt.MapClaims)["expiration"].(string))
	if err != nil {
		return nil, &models.MessageHTTP{HttpCode: 498, Message: "Invalid token"}
	}

	return &TokenData{
		Expiration: expiration,
		UserId:     int(c.Locals(constants.UserId).(float64)),
		Token:      c.Locals(constants.Token).(string),
	}, nil
}

func ValidateToken(tokenString string, c *fiber.Ctx) *models.MessageHTTP {
	c.Locals(constants.Token, tokenString)
	token, err := jwt.Parse(tokenString, validateToken)

	if httpError := normalizeJwtError(token, err); httpError != nil {
		return httpError
	}
	c.Locals(constants.UserId, token.Claims.(jwt.MapClaims)[constants.UserId])
	return nil
}

func normalizeJwtError(token *jwt.Token, err error) *models.MessageHTTP {
	if err != nil {
		return &models.MessageHTTP{HttpCode: fiber.StatusUnauthorized, Message: "Invalid token"}
	}
	switch {
	case token.Valid:
		return nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return &models.MessageHTTP{HttpCode: 498, Message: "That's not even a token"}
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return &models.MessageHTTP{HttpCode: 498, Message: "Invalid token"}
	case errors.Is(err, jwt.ErrTokenExpired):
		return &models.MessageHTTP{HttpCode: fiber.StatusUnauthorized, Message: "Expired"}
	default:
		return &models.MessageHTTP{HttpCode: fiber.StatusInternalServerError, Message: "Unknow error"}
	}
}

type TokenData struct {
	Expiration time.Time
	UserId     int
	Token      string
}
