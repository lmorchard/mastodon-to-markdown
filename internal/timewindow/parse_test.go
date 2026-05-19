package timewindow

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	ref := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)

	got, err := Parse("168h", ref, false)
	if err != nil {
		t.Fatalf("Parse(168h): %v", err)
	}
	want := ref.Add(-168 * time.Hour)
	if !got.Equal(want) {
		t.Errorf("Parse(168h) = %v, want %v", got, want)
	}
}

func TestParseDate(t *testing.T) {
	ref := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)

	got, err := Parse("2026-05-11", ref, false)
	if err != nil {
		t.Fatalf("Parse(2026-05-11): %v", err)
	}
	want := time.Date(2026, 5, 11, 0, 0, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Errorf("Parse(2026-05-11) = %v, want %v", got, want)
	}
}

func TestParseDateEndOfDay(t *testing.T) {
	ref := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)

	got, err := Parse("2026-05-11", ref, true)
	if err != nil {
		t.Fatalf("Parse(2026-05-11, endOfDay): %v", err)
	}
	want := time.Date(2026, 5, 11, 23, 59, 59, int(time.Second-time.Nanosecond), time.Local)
	if !got.Equal(want) {
		t.Errorf("Parse(2026-05-11, endOfDay) = %v, want %v", got, want)
	}
}

func TestParseRFC3339(t *testing.T) {
	ref := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)

	got, err := Parse("2026-05-11T15:30:00Z", ref, false)
	if err != nil {
		t.Fatalf("Parse(RFC3339): %v", err)
	}
	want := time.Date(2026, 5, 11, 15, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Parse(RFC3339) = %v, want %v", got, want)
	}
}

func TestParseEmpty(t *testing.T) {
	if _, err := Parse("", time.Now(), false); err == nil {
		t.Errorf("Parse(\"\") expected error, got nil")
	}
}

func TestParseInvalid(t *testing.T) {
	cases := []string{
		"not-a-date",
		"2026-13-01",   // invalid month
		"2026-05-32",   // invalid day
		"5d",           // not a Go duration ("d" not recognized)
	}
	for _, s := range cases {
		if _, err := Parse(s, time.Now(), false); err == nil {
			t.Errorf("Parse(%q) expected error, got nil", s)
		}
	}
}
