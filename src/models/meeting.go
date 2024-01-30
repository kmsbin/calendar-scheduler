package models

type Meeting struct {
	Id       int                 `json:"id,omitempty"`
	UserId   int                 `json:"user_id,omitempty"`
	Summary  string              `json:"summary"`
	Duration JSONDuration        `json:"duration"`
	Start    JSONTimeRFC3339Nano `json:"start"`
	End      JSONTimeRFC3339Nano `json:"end"`
}
