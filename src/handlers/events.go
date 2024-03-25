package handlers

import (
	"calendar_scheduler/src/constants"
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"calendar_scheduler/src/services"
	"errors"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

const minimunMinutesByRange = 1

func getDatesFromContext(c *fiber.Ctx, meetingsRange *models.MeetingsRange) (*time.Time, *time.Time, error) {
	date, err := time.Parse(time.DateOnly, c.Query("date"))
	if err != nil {
		return nil, nil, models.
			MessageHTTPFromFiberError(fiber.ErrBadRequest).
			FiberContext(c)
	}
	initialTime, finishTime, err := meetingsRange.ConvertToDateRFC3339()
	if err != nil {
		return nil, nil, models.
			MessageHTTPFromFiberError(fiber.ErrBadRequest).
			FiberContext(c)
	}
	*initialTime = setTime(date, *initialTime)
	*finishTime = setTime(date, *finishTime)
	return initialTime, finishTime, nil
}
func setTime(date, timeDate time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		timeDate.Hour(),
		timeDate.Minute(),
		timeDate.Second(),
		0,
		time.Local,
	)
}

func (h Handler) GetMeetingsRangeByContext(c *fiber.Ctx) (*models.MeetingsRange, *models.MessageHTTP) {
	code := c.Query(constants.Code)
	if len(code) == 0 {
		return nil, &models.MessageHTTP{
			HttpCode: fiber.StatusUnprocessableEntity,
			Message:  "missing query parameter",
		}
	}
	meetingsRepository := repositories.NewMeetingsRepository(h.db)
	meetingsRange, err := meetingsRepository.GetmeetingsRangeByCode(code)
	if err != nil {
		if errors.Is(err, repositories.MeetingsRangeNotFounded) {
			return nil, &models.MessageHTTP{
				HttpCode: fiber.StatusNotFound,
				Message:  err.Error(),
			}
		}
		return nil, &models.MessageHTTP{
			HttpCode: fiber.StatusInternalServerError,
			Message:  fiber.ErrInternalServerError.Error(),
		}

	}
	return meetingsRange, nil
}

func (h Handler) GetEventsByCode(c *fiber.Ctx) error {
	meetingsRange, httpModelError := h.GetMeetingsRangeByContext(c)
	if httpModelError != nil {
		return httpModelError.FiberContext(c)
	}
	srv, httpModelError := services.
		NewCalendarServiceFactor(h.db).
		GetCalendarServiceByUserId(meetingsRange.UserId)
	if httpModelError != nil {
		return httpModelError.FiberContext(c)
	}

	initialTime, finishTime, err := getDatesFromContext(c, meetingsRange)
	if err != nil {
		return c.
			Status(fiber.StatusUnprocessableEntity).
			JSON(err)
	}
	events, err := srv.Events.
		List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(initialTime.Format(time.RFC3339)).
		TimeMax(finishTime.Format(time.RFC3339)).
		OrderBy("startTime").
		Do()
	if err != nil {
		log.Printf("Unable to retrieve next ten of the user's events: %v", err)
		return fiber.ErrInternalServerError
	}
	return c.
		Status(200).
		JSON(splitEventsUsingmeetingsRange(getEmptyTimeRange(events.Items, *initialTime, *finishTime), meetingsRange))
}

func revertStartEndOfEvents(events []*calendar.Event) []rangeTimeDate {
	rangeTime := rangeTimeDate{
		models.JSONTimeFromString(events[0].Start.DateTime).Time,
		models.JSONTimeFromString(events[0].End.DateTime).Time,
	}
	countRangeTimes := len(events)
	rangeTimeDates := make([]rangeTimeDate, 0)
	for i := 0; i < countRangeTimes; i++ {
		event := rangeTimeDate{
			models.JSONTimeFromString(events[i].Start.DateTime).Time,
			models.JSONTimeFromString(events[i].End.DateTime).Time,
		}
		if rangeTime.End.After(event.Start) {
			rangeTime.End = event.End
		} else {
			rangeTimeDates = append(rangeTimeDates, rangeTime)
			rangeTime = event
		}
		if countRangeTimes == i+1 {
			rangeTimeDates = append(rangeTimeDates, rangeTime)
		}
	}
	return rangeTimeDates
}

func getEmptyTimeRange(events []*calendar.Event, initialTime, finishTime time.Time) []rangeTimeDate {
	if len(events) == 0 {
		return []rangeTimeDate{{
			Start: initialTime,
			End:   finishTime,
		}}
	}
	rangeTimeDates := revertStartEndOfEvents(events)
	rangeTimeDatesEmpty := make([]rangeTimeDate, 0)
	if rangeTimeDates[0].Start.Sub(initialTime).Minutes() > minimunMinutesByRange {
		rangeTimeDatesEmpty = append(rangeTimeDatesEmpty, rangeTimeDate{
			Start: initialTime,
			End:   rangeTimeDates[0].Start,
		})
	}

	for i := 0; i < len(rangeTimeDates)-1; i++ {
		currentTimeRange := rangeTimeDates[i]
		rangeTimeDatesEmpty = append(rangeTimeDatesEmpty, rangeTimeDate{
			Start: currentTimeRange.End,
			End:   rangeTimeDates[i+1].Start,
		})
	}
	if len(rangeTimeDates) == 0 {
		return rangeTimeDates
	}
	lastRange := rangeTimeDates[len(rangeTimeDates)-1]
	if finishTime.Sub(lastRange.End).Minutes() > minimunMinutesByRange {
		rangeTimeDatesEmpty = append(rangeTimeDatesEmpty, rangeTimeDate{
			Start: lastRange.End,
			End:   finishTime,
		})
	}

	return rangeTimeDatesEmpty
}

func splitEventsUsingmeetingsRange(events []rangeTimeDate, meetingsRange *models.MeetingsRange) []rangeTimeDate {
	meetingsDuration := meetingsRange.Duration.Duration()
	splittedRanges := make([]rangeTimeDate, 0)

	for _, rangeEvent := range events {
		for rangeEvent.End.Sub(rangeEvent.Start).Minutes() > meetingsDuration.Minutes() {
			endDurationDate := rangeEvent.Start.Add(meetingsDuration)
			splittedRanges = append(splittedRanges, rangeTimeDate{
				Start: rangeEvent.Start,
				End:   endDurationDate,
			})
			rangeEvent.Start = endDurationDate
		}
	}

	return splittedRanges
}

type rangeTimeDate struct {
	Start time.Time `json:"start_date"`
	End   time.Time `json:"end_date"`
}
