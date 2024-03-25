package models

import "time"

type ResetPassword struct {
	UserId int       `json:"user_id"`
	Email  string    `json:"email"`
	Code   string    `json:"code"`
	Expiry time.Time `json:"expiry"`
}
