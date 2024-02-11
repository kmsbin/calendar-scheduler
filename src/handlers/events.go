package handlers

import (
	"calendar_scheduler/src/models"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

const minimunMinutesByRange = 1

func getDatesFromContext(c *fiber.Ctx) (*time.Time, *time.Time, error) {
	initialTime, err := time.Parse(time.RFC3339, c.Query("initial_date"))
	if err != nil {
		return nil, nil, models.MessageHTTPFromFiberError(fiber.ErrBadRequest)
	}
	finishTime, err := time.Parse(time.RFC3339, c.Query("finish_date"))
	if err != nil {
		return &initialTime, nil, models.MessageHTTPFromFiberError(fiber.ErrBadRequest)
	}
	return &initialTime, &finishTime, nil
}

func GetEmptyScheduledTime(c *fiber.Ctx) error {
	srv, httpModelError := GetCalendarService(c)
	if httpModelError != nil {
		return c.Status(httpModelError.HttpCode).JSON(httpModelError)
	}
	initialTime, finishTime, err := getDatesFromContext(c)
	if err != nil {
		return nil
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
	if len(events.Items) == 0 {
		return c.Status(200).JSON([]rangeTimeDate{})
	}
	return c.
		Status(200).
		JSON(getEmptyTimeRange(events.Items, *initialTime, *finishTime))
}

func splitEvents(events []*calendar.Event) []rangeTimeDate {
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
	rangeTimeDates := splitEvents(events)
	rangeTimeDatesEmpty := make([]rangeTimeDate, 0)
	if rangeTimeDates[0].Start.Sub(initialTime).Minutes() > minimunMinutesByRange {
		rangeTimeDatesEmpty = append(rangeTimeDatesEmpty, rangeTimeDate{
			Start: initialTime,
			End:   rangeTimeDates[0].Start,
		})
	}

	for i := 0; i < len(rangeTimeDates); i++ {
		currentTimeRange := rangeTimeDates[i]
		if len(rangeTimeDates) == i+1 {
			if finishTime.Sub(currentTimeRange.End).Minutes() > minimunMinutesByRange {
				rangeTimeDatesEmpty = append(rangeTimeDatesEmpty, rangeTimeDate{
					Start: currentTimeRange.End,
					End:   finishTime,
				})
			}
		} else {
			rangeTimeDatesEmpty = append(rangeTimeDatesEmpty, rangeTimeDate{
				Start: currentTimeRange.End,
				End:   rangeTimeDates[i+1].Start,
			})
		}
	}

	return rangeTimeDatesEmpty
}

type rangeTimeDate struct {
	Start time.Time `json:"start_date"`
	End   time.Time `json:"end_date"`
}
