package handlers

import (
	"calendar_scheduler/src/models"
	"calendar_scheduler/src/repositories"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

const minimunMinutesByRange = 1

func getDatesFromContext(c *fiber.Ctx, meetingRange *models.MeetingRange) (*time.Time, *time.Time, error) {
	date, err := time.Parse(time.DateOnly, c.Query("date"))
	if err != nil {
		log.Printf("query error, param: %v, error: %v", c.Query("date"), err)
		return nil, nil, models.MessageHTTPFromFiberError(fiber.ErrBadRequest)
	}
	initialTime, finishTime, err := meetingRange.ConvertToDateRFC3339()
	if err != nil {
		log.Printf("meeting range, start: %v, end: %v, error: %v", meetingRange.Start, meetingRange.End, err)
		return nil, nil, models.MessageHTTPFromFiberError(fiber.ErrBadRequest)
	}
	log.Printf("duration %v", meetingRange.Duration.Duration().Minutes())
	*initialTime = setTime(date, *initialTime)
	*finishTime = setTime(date, *finishTime)
	log.Printf("initialTime %v", initialTime)
	log.Printf("finishTime %v", finishTime)
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

func GetEmptyScheduledTime(c *fiber.Ctx) error {
	srv, httpModelError := GetCalendarService(c)
	userId := c.Locals("user_id")
	if httpModelError != nil {
		return c.Status(httpModelError.HttpCode).JSON(httpModelError)
	}
	meetingRepository := repositories.NewMeetingRepository()
	meetingRange, err := meetingRepository.GetLastMeetingRange(userId)
	if err != nil {
		return c.
			Status(fiber.StatusPreconditionRequired).
			JSON(models.MessageHTTP{
				Message: "is needed create a meeting range before",
			})
	}
	initialTime, finishTime, err := getDatesFromContext(c, meetingRange)
	if err != nil {
		log.Printf("date format error %v", err)
		return c.JSON(err)
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
	//if len(events.Items) == 0 {
	//	return c.Status(200).JSON([]rangeTimeDate{})
	//}
	return c.
		Status(200).
		JSON(splitEventsUsingMeetingRange(getEmptyTimeRange(events.Items, *initialTime, *finishTime), meetingRange))
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

func splitEventsUsingMeetingRange(events []rangeTimeDate, meetingRange *models.MeetingRange) []rangeTimeDate {
	meetingDuration := meetingRange.Duration.Duration()
	splittedRanges := make([]rangeTimeDate, 0)

	for _, rangeEvent := range events {
		for rangeEvent.End.Sub(rangeEvent.Start).Minutes() > meetingDuration.Minutes() {
			endDurationDate := rangeEvent.Start.Add(meetingDuration)
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
