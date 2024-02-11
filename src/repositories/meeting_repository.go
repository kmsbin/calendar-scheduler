package repositories

import (
	"calendar_scheduler/src/database"
	"calendar_scheduler/src/models"
	"database/sql"
	"errors"
	"log"
	"time"
)

type MeetingRepository struct {
	db *sql.DB
}

func NewMeetingRepository() MeetingRepository {
	db, err := database.OpenConnection()
	if err != nil {
		panic(err)
	}
	return MeetingRepository{db}
}

func (m *MeetingRepository) InsertMeetingRange(meetingBody models.MeetingRange) error {
	_, err := m.db.Exec(
		"insert into meeting_range(user_id, summary, start_time, end_time, duration) values ($1, $2, $3, $4, $5)",
		meetingBody.UserId,
		meetingBody.Summary,
		meetingBody.Start,
		meetingBody.End,
		meetingBody.Duration,
	)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (m *MeetingRepository) GetLastMeetingRange(userId any) (*models.MeetingRange, error) {
	row := m.db.QueryRow("select id, user_id, summary, start_time, end_time, duration from meeting_range where user_id = $1", userId)

	meetingRange := models.MeetingRange{}
	var duration float64
	err := row.Scan(
		&meetingRange.Id,
		&meetingRange.UserId,
		&meetingRange.Summary,
		&meetingRange.Start,
		&meetingRange.End,
		&duration,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	meetingRange.Duration = models.JSONDuration(time.Duration(duration))
	return &meetingRange, nil
}
