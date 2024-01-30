package auth

import (
	"calendar_scheduler/src/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	SigningKey = "$argon2i$v=19$m=16,t=2,p=2$a2F1bGluZG8$wgUw6YD/lSRkRy23XQ/JH6hugjyqCbm1"
)

type UserClaims struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func CreateToken(user models.User) (string, error) {
	claims := UserClaims{
		UserId: user.Id,
		Email:  user.Email,
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 6))
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString([]byte(SigningKey))
	if err != nil {
		return "", err
	}
	return signedString, nil
}
