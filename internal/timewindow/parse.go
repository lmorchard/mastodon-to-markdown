// Package timewindow parses the canonical --since/--until flag values used by
// the orchestrator-facing `export` subcommand.
//
// This package is intentionally duplicated across the *-to-markdown tools
// (one copy per repo) rather than shared via module — the contract is small,
// stable, and not worth a cross-repo dependency.
package timewindow

import (
	"fmt"
	"time"
)

// Parse interprets s as one of:
//   - a Go duration string ("168h", "30m") — returned as ref.Add(-d)
//     (durations are interpreted as "ago" relative to ref);
//   - an RFC3339 timestamp;
//   - a YYYY-MM-DD date in the local timezone.
//
// When endOfDay is true, a YYYY-MM-DD date is advanced to 23:59:59.999999999
// of that day. RFC3339 and durations ignore endOfDay.
//
// Empty s returns an error; callers that want a default should check for ""
// before calling Parse.
func Parse(s string, ref time.Time, endOfDay bool) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("empty time value")
	}
	if d, err := time.ParseDuration(s); err == nil {
		return ref.Add(-d), nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		if endOfDay {
			t = t.Add(24*time.Hour - time.Nanosecond)
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid time %q (expected YYYY-MM-DD, RFC3339, or Go duration)", s)
}
