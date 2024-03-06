package models

type CalendarEvent struct {
	Email string   `json:"email"`
	Date  JSONTime `json:"date"`
}
