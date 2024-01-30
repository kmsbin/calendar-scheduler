package models

type CalendarEvent struct {
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Email       string   `json:"email"`
	Start       JSONTime `json:"start"`
	End         JSONTime `json:"end"`
}
