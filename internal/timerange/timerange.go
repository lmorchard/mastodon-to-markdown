package timerange

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// TimeRange represents a time period with start and end times
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Parse creates a TimeRange from the given parameters
// Priority: if start/end are provided, use them; otherwise use since duration
func Parse(since, start, end string) (*TimeRange, error) {
	now := time.Now()

	// If start and end are both provided, use them
	if start != "" && end != "" {
		startTime, err := parseDate(start)
		if err != nil {
			return nil, fmt.Errorf("invalid start date: %w", err)
		}
		endTime, err := parseDate(end)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
		if endTime.Before(startTime) {
			return nil, fmt.Errorf("end date must be after start date")
		}
		return &TimeRange{Start: startTime, End: endTime}, nil
	}

	// If only start is provided, end defaults to now
	if start != "" {
		startTime, err := parseDate(start)
		if err != nil {
			return nil, fmt.Errorf("invalid start date: %w", err)
		}
		return &TimeRange{Start: startTime, End: now}, nil
	}

	// If since is provided, calculate start from duration
	if since != "" {
		duration, err := parseDuration(since)
		if err != nil {
			return nil, fmt.Errorf("invalid since duration: %w", err)
		}
		return &TimeRange{Start: now.Add(-duration), End: now}, nil
	}

	// Default: last 7 days
	return &TimeRange{Start: now.AddDate(0, 0, -7), End: now}, nil
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("expected format YYYY-MM-DD, got: %s", dateStr)
	}
	return t, nil
}

// parseDuration parses duration strings like "24h", "7d", "2w"
// Supports: h (hours), d (days), w (weeks)
func parseDuration(s string) (time.Duration, error) {
	re := regexp.MustCompile(`^(\d+)([hdw])$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid duration format: expected format like '24h', '7d', or '2w'")
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %w", err)
	}

	unit := matches[2]
	switch unit {
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	case "w":
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid duration unit: %s", unit)
	}
}

// FormatDate formats a time as YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime formats a time with date and time
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}
