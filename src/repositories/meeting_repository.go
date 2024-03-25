package repositories

import (
	"calendar_scheduler/src/database"
	"calendar_scheduler/src/models"
	"database/sql"
	"errors"
	"log"
)

type MeetingsRepository struct {
	db *sql.DB
}

func NewMeetingsRepository(db *sql.DB) MeetingsRepository {
	return MeetingsRepository{db}
}

func (m *MeetingsRepository) InsertMeetingsRange(meetingsBody models.MeetingsRange) error {
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

func (m *MeetingsRepository) InsertMeetingsRangeEmail(meetingsBody models.MeetingRangeEmail) error {
	_, err := m.db.Exec(
		"insert into meetings_ranges_emails(user_id, meetings_id, email) values ($1, $2, $3)",
		meetingsBody.UserId,
		meetingsBody.MeetingId,
		meetingsBody.Email,
	)
	if database.IsPqErrorCode(err, database.UniqueViolationErr) {
		return EmailDuplicatedInMeetingRange
	}

	return err
}

func (m *MeetingsRepository) GetLastmeetingsRange(userId any) (*models.MeetingsRange, error) {
	row := m.db.QueryRow("select id, code, user_id, summary, start_time, end_time, duration from meetings_ranges where user_id = $1", userId)
	return scanTomeetingsRange(row)
}

func scanTomeetingsRange(row *sql.Row) (*models.MeetingsRange, error) {
	meetingsRange := models.MeetingsRange{}
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
			return nil, MeetingsRangeNotFounded
		}
		return nil, err
	}
	if &meetingsRange == nil {
		log.Printf("meetingsRange == nil %v", meetingsRange)
		return nil, MeetingsRangeNotFounded
	}
	log.Printf("meetingsRange %v", meetingsRange)
	meetingsRange.Duration = models.JSONDuration(duration)
	return &meetingsRange, nil
}

func (m *MeetingsRepository) GetmeetingsRangeByCode(code string) (*models.MeetingsRange, error) {
	row := m.db.QueryRow("select id, code, user_id, summary, start_time, end_time, duration from meetings_ranges where code = $1", code)
	return scanTomeetingsRange(row)
}
