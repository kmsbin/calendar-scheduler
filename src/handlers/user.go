package handlers

import (
	"calendar_scheduler/src/auth"
	"calendar_scheduler/src/database"
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
)

func CreateUserHadler(c *fiber.Ctx) error {
	c.Set(headers.ContentType, fiber.MIMEApplicationJSON)
	var user models.User
	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		log.Print(err)
		return fiber.ErrNotAcceptable
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	db, _ := database.OpenConnection()
	_, err = db.Exec("insert into users(name, email, password) values ($1, $2, $3)", user.Name, user.Email, password)
	if err != nil {
		log.Print(err)
		return fiber.ErrBadGateway
	}
	return c.JSON(map[string]string{"message": "Successful!!"})
}

func GetUser(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id")
	if userId == nil {
		return fiber.ErrUnauthorized
	}
	userRepository := repositories.NewUserRepository()
	user, err := userRepository.GetUserById(userId.(int))
	if err != nil {
		return ctx.
			Status(fiber.StatusNotFound).
			JSON(models.MessageHTTPFromFiberError(fiber.ErrNotFound))
	}
	return ctx.Status(fiber.StatusOK).JSON(user)
}
func SignInUser(c *fiber.Ctx) error {
	email, password := c.Query("email"), c.Query("password")
	if email == "" || password == "" {
		return fiber.ErrNotAcceptable
	}
	userRepository := repositories.NewUserRepository()
	user, userPassword, err := userRepository.GetUserByEmail(email)
	if err != nil {
		log.Printf("error %v\n", err)
		if errors.Is(err, repositories.UserNotFounded) {
			return c.
				Status(fiber.StatusNotFound).
				JSON(models.MessageHTTPFromFiberError(fiber.ErrNotFound))
		}
		return models.MessageHTTPFromFiberError(fiber.ErrInternalServerError)
	}
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	if err != nil {
		return models.MessageHTTPFromFiberError(fiber.ErrUnauthorized)
	}
	token, err := auth.CreateToken(*user)
	if err != nil {
		return models.MessageHTTPFromFiberError(fiber.ErrUnauthorized)
	}
	return c.JSON(map[string]string{"token": token})
}

func ValidateTokenMiddleware(c *fiber.Ctx) error {
	tokenString, ok := c.Queries()["token"]
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(map[string]string{"message": "Forbidden!"})
	}
	return ValidateToken(tokenString, c)
}

func ValidateToken(tokenString string, c *fiber.Ctx) error {
	c.Locals("token", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.SigningKey), nil
	})
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"message": "Invalid token"})
	}
	switch {
	case token.Valid:
		c.Locals("user_id", token.Claims.(jwt.MapClaims)["user_id"])
		return c.Next()
	case errors.Is(err, jwt.ErrTokenMalformed):
		return c.Status(498).JSON(map[string]string{"message": "That's not even a token"})
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return c.Status(498).JSON(map[string]string{"message": "Invalid token"})
	case errors.Is(err, jwt.ErrTokenExpired):
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"message": "Expired"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"message": "Unknow error "})
	}
}
