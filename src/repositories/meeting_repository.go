package repositories

import (
	"calendar_scheduler/src/models"
	"database/sql"
	"errors"
	"log"
)

type meetingsRepository struct {
	db *sql.DB
}

func NewmeetingsRepository(db *sql.DB) meetingsRepository {
	return meetingsRepository{db}
}

func (m *meetingsRepository) InsertmeetingsRange(meetingsBody models.meetingsRange) error {
	_, err := m.db.Exec(
		"insert into meetings_ranges(user_id, summary, start_time, end_time, duration, code) values ($1, $2, $3, $4, $5, $6)",
		meetingsBody.UserId,
		meetingsBody.Summary,
		meetingsBody.Start,
		meetingsBody.End,
		meetingsBody.Duration,
		meetingsBody.Code,
	)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (m *meetingsRepository) GetLastmeetingsRange(userId any) (*models.meetingsRange, error) {
	row := m.db.QueryRow("select id, code, user_id, summary, start_time, end_time, duration from meetings_ranges where user_id = $1", userId)
	return scanTomeetingsRange(row)
}

func scanTomeetingsRange(row *sql.Row) (*models.meetingsRange, error) {
	meetingsRange := models.meetingsRange{}
	var duration float64
	err := row.Scan(
		&meetingsRange.Id,
		&meetingsRange.Code,
		&meetingsRange.UserId,
		&meetingsRange.Summary,
		&meetingsRange.Start,
		&meetingsRange.End,
		&duration,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, meetingsRangeNotFounded
		}
		return nil, err
	}
	if &meetingsRange == nil {
		log.Printf("meetingsRange == nil %v", meetingsRange)
		return nil, meetingsRangeNotFounded
	}
	log.Printf("meetingsRange %v", meetingsRange)
	meetingsRange.Duration = models.JSONDuration(duration)
	return &meetingsRange, nil
}

func (m *meetingsRepository) GetmeetingsRangeByCode(code string) (*models.meetingsRange, error) {
	row := m.db.QueryRow("select id, code, user_id, summary, start_time, end_time, duration from meetings_ranges where code = $1", code)
	return scanTomeetingsRange(row)
}
