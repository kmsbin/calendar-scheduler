package handlers

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"calendar_scheduler/src/services"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

func (h Handler) CreateGoogleCalendarEvent(c *fiber.Ctx) error {
	calendarEvent := models.CalendarEvent{}
	if err := c.BodyParser(&calendarEvent); err != nil {
		return UnprocessableEntity(c)
	}
	meetingsRange, httpModelError := h.GetMeetingsRangeByContext(c)
	if httpModelError != nil {
		return httpModelError.FiberContext(c)
	}

	meetingRepository := repositories.NewMeetingsRepository(h.db)

	err := meetingRepository.InsertMeetingsRangeEmail(models.MeetingRangeEmail{
		UserId:    meetingsRange.UserId,
		MeetingId: meetingsRange.Id,
		Email:     calendarEvent.Email,
	})

	if err != nil {
		if errors.Is(err, repositories.EmailDuplicatedInMeetingRange) {
			return ConflictError(c, models.MessageHTTPFromMessage(err.Error()))
		}
		return InternalServerError(c)
	}

	srv, httpModelError := services.
		NewCalendarServiceFactor(h.db, c.BaseURL()).
		GetCalendarServiceByUserId(meetingsRange.UserId)
	if httpModelError != nil {
		return httpModelError.FiberContext(c)
	}

	event := prepareCalendarEvent(calendarEvent, meetingsRange)
	fmt.Printf("Event created: %s\n", event.HangoutLink)
	event, err = srv.Events.
		Insert(constants.CalendarId, event).
		ConferenceDataVersion(1).
		SendUpdates("all").
		Do()
	if err != nil {
		log.Printf("Unable to create event. %v\n", err)
		return InternalServerError(c)
	}
	fmt.Printf("Event created: %s\n", event.HangoutLink)
	return ResponseOK(c, models.MessageHTTPFromMessage(event.HangoutLink))
}

func prepareCalendarEvent(calendarEvent models.CalendarEvent, meetingsRange *models.MeetingsRange) *calendar.Event {
	meetingDuration, err := time.ParseDuration(meetingsRange.Duration)
	if err != nil {
		panic(err)
	}
	return &calendar.Event{
		Summary: meetingsRange.Summary,
		Start: &calendar.EventDateTime{
			DateTime: calendarEvent.Date.Format(time.RFC3339),
			TimeZone: constants.Locale,
		},
		End: &calendar.EventDateTime{
			DateTime: calendarEvent.Date.
				Add(meetingDuration).
				Format(time.RFC3339), //  calendarEvent.End.Format(time.RFC3339),
			TimeZone: constants.Locale,
		},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
				RequestId: uuid.New().String(),
			},
		},
		Attendees: []*calendar.EventAttendee{
			{Email: calendarEvent.Email},
		},
	}
}
