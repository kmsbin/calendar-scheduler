package models

type MeetingRange struct {
	Id       int          `json:"id,omitempty"`
	UserId   int          `json:"user_id,omitempty"`
	Summary  string       `json:"summary"`
	Duration JSONDuration `json:"duration"`
	Start    string       `json:"start"`
	End      string       `json:"end"`
}
