package utils

import (
	"errors"
	"time"
)

// ParseDate parses date string in various formats
func ParseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",           // YYYY-MM-DD
		"2006-01-02T15:04:05Z", // RFC3339
		"02/01/2006",           // DD/MM/YYYY
		"01/02/2006",           // MM/DD/YYYY
		"2006-01-02 15:04:05",  // YYYY-MM-DD HH:MM:SS
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("invalid date format")
}

// ParseDateTime parses datetime string in RFC3339 format
func ParseDateTime(datetimeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, datetimeStr)
}

// FormatDate formats time to date string
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime formats time to datetime string
func FormatDateTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// StartOfMonth returns the first day of the month
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the last day of the month
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Second)
}

// GetCurrentMonth returns start and end of current month
func GetCurrentMonth() (time.Time, time.Time) {
	now := time.Now()
	start := StartOfMonth(now)
	end := EndOfMonth(now)
	return start, end
}
