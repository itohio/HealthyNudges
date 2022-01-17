package nudge

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	weekDays = map[time.Weekday]string{
		time.Sunday:    strings.ToLower(time.Sunday.String()),
		time.Monday:    strings.ToLower(time.Monday.String()),
		time.Tuesday:   strings.ToLower(time.Tuesday.String()),
		time.Wednesday: strings.ToLower(time.Wednesday.String()),
		time.Thursday:  strings.ToLower(time.Thursday.String()),
		time.Friday:    strings.ToLower(time.Friday.String()),
		time.Saturday:  strings.ToLower(time.Saturday.String()),
	}
	scheduleParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
)

// Greedy fuzzy search
func expandWeekDay(dayNames map[time.Weekday]string, shorthand string) (time.Weekday, error) {
	shorthand = strings.ToLower(shorthand)

	for i, day := range weekDays {
		if strings.HasPrefix(day, shorthand) {
			return time.Weekday(i), nil
		}
	}

	return 0, fmt.Errorf("week day not matched")
}

func ValidateSchedule(schedule string) error {
	for _, part := range strings.Split(schedule, ";") {
		if _, err := scheduleParser.Parse(part); err != nil {
			return err
		}
	}
	return nil
}

// Parse and match cron schedule definition
// ┌───────────── minute (0 - 59)
// │ ┌───────────── hour (0 - 23)
// │ │ ┌───────────── day of the month (1 - 31)
// │ │ │ ┌───────────── month (1 - 12)
// │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday;
// │ │ │ │ │                                   7 is also Sunday on some systems)
// │ │ │ │ │
// │ │ │ │ │
// * * * * *
// Each value may specify an interval, e.g.:
// * 18-08 * * 1-5
// Translates to this: from 18 till 8 from Monday till Friday
func MatchSchedule(schedule string, advance time.Duration) bool {
	for _, part := range strings.Split(schedule, ";") {
		sc, err := scheduleParser.Parse(part)
		if err != nil {
			return false
		}

		now := time.Now().Add(advance)
		next := sc.Next(now.Add(-time.Minute))
		if next.Sub(now).Minutes() <= 0 {
			return true
		}
	}
	return false
}
