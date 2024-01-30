package models

import (
	"errors"
	"time"
)

type JSONTime struct {
	time.Time
}

type JSONTimeRFC3339Nano struct {
	JSONTime
}

func JSONTimeFromString(value string) JSONTime {
	parsedTime, _ := time.Parse(time.RFC3339, value)
	return JSONTime{parsedTime}
}

func (t *JSONTime) ToRFC3339Nano() JSONTimeRFC3339Nano {
	return JSONTimeRFC3339Nano{*t}
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	err := t.unmarshalTimeLayout(time.RFC3339, data)
	if err == nil {
		return nil
	}
	return t.unmarshalTimeLayout(time.RFC3339Nano, data)
}

func (t *JSONTimeRFC3339Nano) UnmarshalJSON(data []byte) error {
	err := t.unmarshalTimeLayout(time.RFC3339Nano, data)
	if err == nil {
		return nil
	}
	return t.unmarshalTimeLayout(time.RFC3339, data)
}

func (t *JSONTime) unmarshalTimeLayout(layout string, data []byte) error {
	err := t.Time.UnmarshalJSON(data)
	if err == nil {
		return nil
	}
	if string(data) == "null" {
		return nil
	}
	ti, err := time.Parse(`"`+layout+`"`, string(data))

	*t = JSONTime{ti}
	return err
}

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	bytes, err := t.Time.MarshalJSON()
	if err == nil {
		return bytes, err
	}
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(time.RFC3339)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.RFC3339)
	b = append(b, '"')
	return b, nil
}

func (t *JSONTimeRFC3339Nano) ToString() string {
	return t.Format(time.RFC3339Nano)
}

func (t *JSONTime) ToString() string {
	return t.Format(time.RFC3339)
}
