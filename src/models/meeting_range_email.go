package models

type MeetingRangeEmail struct {
	MeetingId int    `json:"meetings_id,omitempty"`
	UserId    int    `json:"user_id,omitempty"`
	Email     string `json:"email"`
}
