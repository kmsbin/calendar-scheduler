package models

import "time"

type MeetingRange struct {
	Id       int          `json:"id,omitempty"`
	UserId   int          `json:"user_id,omitempty"`
	Summary  string       `json:"summary"`
	Duration JSONDuration `json:"duration"`
	Start    string       `json:"start"`
	End      string       `json:"end"`
}

func (m *MeetingRange) ConvertToDateRFC3339() (*time.Time, *time.Time, error) {
	start, err := time.Parse(time.RFC3339, m.Start)
	if err != nil {
		return nil, nil, err
	}
	end, err := time.Parse(time.RFC3339, m.End)
	if err != nil {
		return nil, nil, err
	}
	return &start, &end, nil
}

func (m *MeetingRange) ConvertToTime() (*time.Time, *time.Time, error) {
	start, err := time.Parse(time.TimeOnly, m.Start)
	if err != nil {
		return nil, nil, err
	}
	end, err := time.Parse(time.TimeOnly, m.End)
	if err != nil {
		return nil, nil, err
	}
	return &start, &end, nil
}
