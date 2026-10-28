package tides

import (
	"testing"
	"time"
)

func TestParseOpenMeteoUTCEventsUsesCanaryDST(t *testing.T) {
	canary, err := time.LoadLocation("Atlantic/Canary")
	if err != nil {
		t.Fatalf("load Atlantic/Canary: %v", err)
	}

	tests := []struct {
		name           string
		date           time.Time
		utcTimes       []string
		expectedClock  string
		expectedOffset int
	}{
		{
			name:           "summer time before transition",
			date:           time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC),
			utcTimes:       []string{"2026-10-24T11:00", "2026-10-24T11:15", "2026-10-24T11:30"},
			expectedClock:  "2026-10-24 12:15",
			expectedOffset: 60 * 60,
		},
		{
			name:           "standard time after transition",
			date:           time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC),
			utcTimes:       []string{"2026-10-25T11:30", "2026-10-25T11:45", "2026-10-25T12:00"},
			expectedClock:  "2026-10-25 11:45",
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startLocal := time.Date(tt.date.Year(), tt.date.Month(), tt.date.Day(), 0, 0, 0, 0, canary)
			endLocal := startLocal.AddDate(0, 0, 1)
			payload := openMeteoPayload{}
			payload.Minutely15.Time = tt.utcTimes
			payload.Minutely15.Height = []float64{0.7, 0.8, 0.7}

			events, err := parseOpenMeteoUTCEvents(payload, startLocal, endLocal)
			if err != nil {
				t.Fatalf("parse UTC events: %v", err)
			}
			if len(events) != 1 {
				t.Fatalf("expected one event, got %d", len(events))
			}
			if got := events[0].TimeLocal.Format("2006-01-02 15:04"); got != tt.expectedClock {
				t.Fatalf("expected local time %s, got %s", tt.expectedClock, got)
			}
			_, offset := events[0].TimeLocal.Zone()
			if offset != tt.expectedOffset {
				t.Fatalf("expected UTC offset %d, got %d", tt.expectedOffset, offset)
			}
		})
	}
}

func TestCanaryDSTEndDayHas25Hours(t *testing.T) {
	canary, err := time.LoadLocation("Atlantic/Canary")
	if err != nil {
		t.Fatalf("load Atlantic/Canary: %v", err)
	}

	startLocal := time.Date(2026, 10, 25, 0, 0, 0, 0, canary)
	endLocal := startLocal.AddDate(0, 0, 1)
	if got := endLocal.Sub(startLocal); got != 25*time.Hour {
		t.Fatalf("expected DST transition day to last 25 hours, got %s", got)
	}
	if got := startLocal.UTC().Format("2006-01-02 15:04"); got != "2026-10-24 23:00" {
		t.Fatalf("unexpected UTC start: %s", got)
	}
	if got := endLocal.UTC().Format("2006-01-02 15:04"); got != "2026-10-26 00:00" {
		t.Fatalf("unexpected UTC end: %s", got)
	}
}
