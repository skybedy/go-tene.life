package web

import (
	"testing"
	"time"

	"github.com/skybedy/laravel-tene.life/internal/models"
)

func TestBuildHomepageTideEventViewsLimitsEachTypeToTwoEvents(t *testing.T) {
	eventTime := func(hour, minute int) time.Time {
		return time.Date(2026, time.November, 18, hour, minute, 0, 0, time.FixedZone("WET", 0))
	}
	events := []models.TideEvent{
		{EventType: "HIGH", EventTimeLocal: eventTime(1, 5), HeightM: 1.21},
		{EventType: "LOW", EventTimeLocal: eventTime(4, 10), HeightM: -0.34},
		{EventType: "HIGH", EventTimeLocal: eventTime(8, 15), HeightM: 1.43},
		{EventType: "LOW", EventTimeLocal: eventTime(11, 20), HeightM: -0.28},
		{EventType: "HIGH", EventTimeLocal: eventTime(15, 25), HeightM: 1.36},
		{EventType: "LOW", EventTimeLocal: eventTime(18, 30), HeightM: -0.31},
	}

	high, low := buildHomepageTideEventViews(events)

	if len(high) != maxHomepageTideEventsPerType {
		t.Fatalf("high tide count = %d, want %d", len(high), maxHomepageTideEventsPerType)
	}
	if len(low) != maxHomepageTideEventsPerType {
		t.Fatalf("low tide count = %d, want %d", len(low), maxHomepageTideEventsPerType)
	}
	if high[0].Time != "1:05" || high[1].Time != "8:15" {
		t.Fatalf("high tide times = %#v, want first two events", high)
	}
	if low[0].Time != "4:10" || low[1].Time != "11:20" {
		t.Fatalf("low tide times = %#v, want first two events", low)
	}
}
