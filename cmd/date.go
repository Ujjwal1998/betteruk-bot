package cmd

import (
	"fmt"
	"os"
	"time"
)

const maxBookingDaysAhead = 5

func allowedBookingDates(now time.Time) []string {
	dates := make([]string, maxBookingDaysAhead+1)
	for i := 0; i <= maxBookingDaysAhead; i++ {
		dates[i] = now.AddDate(0, 0, i).Format("2006-01-02")
	}
	return dates
}

func validateBookingDate(date string, now time.Time) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("invalid date %q: use YYYY-MM-DD", date)
	}
	for _, allowed := range allowedBookingDates(now) {
		if date == allowed {
			return nil
		}
	}
	return fmt.Errorf("date must be today or up to %d days ahead (use -d %s … %s)",
		maxBookingDaysAhead,
		allowedBookingDates(now)[0],
		allowedBookingDates(now)[maxBookingDaysAhead],
	)
}

func printCurrentDate(date string) {
	if date == "" {
		return
	}
	fmt.Fprintf(os.Stderr, "Date: %s  [%s]\n\n", formatDateLabel(date, time.Now()), date)
}

func formatDateLabel(date string, now time.Time) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	label := t.Format("Mon 2 Jan 2006")
	switch date {
	case now.Format("2006-01-02"):
		return label + " (today)"
	case now.AddDate(0, 0, 1).Format("2006-01-02"):
		return label + " (tomorrow)"
	default:
		return label
	}
}
