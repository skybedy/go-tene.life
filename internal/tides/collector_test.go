package tides

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/skybedy/laravel-tene.life/internal/models"
)

type fakeRepo struct {
	events      []models.TideEvent
	hasFresh    bool
	replaceRuns int
}

func (f *fakeRepo) GetTideEvents(ctx context.Context, dateLocal, locationKey string) ([]models.TideEvent, error) {
	return f.events, nil
}

func (f *fakeRepo) HasFreshPuertosData(ctx context.Context, dateLocal, locationKey string, sinceUTC time.Time) (bool, error) {
	return f.hasFresh, nil
}

func (f *fakeRepo) ReplaceTideEvents(ctx context.Context, dateLocal, locationKey string, events []models.TideEvent) error {
	f.replaceRuns++
	f.events = append([]models.TideEvent(nil), events...)
	return nil
}

func TestCollectTidesFallbackToOpenMeteo(t *testing.T) {
	repo := &fakeRepo{}
	collector := NewCollector(repo)
	collector.SetServingSource("hybrid")

	madrid, _ := time.LoadLocation("Europe/Madrid")
	collector.SetFetchers(
		func(ctx context.Context, date time.Time, loc LocationConfig) ([]ExtremeEvent, string, error) {
			return nil, "", errors.New("puertos down")
		},
		func(ctx context.Context, date time.Time, loc LocationConfig) ([]ExtremeEvent, string, error) {
			return []ExtremeEvent{
				{Type: "LOW", TimeLocal: time.Date(2026, 3, 1, 6, 0, 0, 0, madrid), HeightM: -1.1},
				{Type: "HIGH", TimeLocal: time.Date(2026, 3, 1, 12, 0, 0, 0, madrid), HeightM: 0.5},
			}, `{"source":"open"}`, nil
		},
	)

	if err := collector.CollectTides(context.Background(), "2026-03-01", "los_cristianos"); err != nil {
		t.Fatalf("CollectTides failed: %v", err)
	}
	if repo.replaceRuns != 1 {
		t.Fatalf("expected one replace run, got %d", repo.replaceRuns)
	}
	if len(repo.events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(repo.events))
	}
	if repo.events[0].Source != "open_meteo" {
		t.Fatalf("expected open_meteo source, got %s", repo.events[0].Source)
	}
	if repo.events[0].Confidence != 70 {
		t.Fatalf("expected confidence 70, got %d", repo.events[0].Confidence)
	}
}

func TestValidateEventsForDateRejectsOutsideDay(t *testing.T) {
	madrid, _ := time.LoadLocation("Europe/Madrid")
	events := []ExtremeEvent{
		{Type: "LOW", TimeLocal: time.Date(2026, 3, 1, 6, 0, 0, 0, madrid), HeightM: -1.1},
		{Type: "HIGH", TimeLocal: time.Date(2026, 3, 2, 0, 5, 0, 0, madrid), HeightM: 0.5},
	}

	err := validateEventsForDate("2026-03-01", "Europe/Madrid", events)
	if err == nil {
		t.Fatal("expected validation error for out-of-date event")
	}
}

func TestCollectTidesPuertosOverwritesOpenMeteo(t *testing.T) {
	madrid, _ := time.LoadLocation("Europe/Madrid")
	repo := &fakeRepo{
		events: []models.TideEvent{
			{DateLocal: "2026-03-01", LocationKey: "los_cristianos", EventType: "LOW", Source: "open_meteo", EventTimeLocal: time.Date(2026, 3, 1, 6, 0, 0, 0, madrid)},
			{DateLocal: "2026-03-01", LocationKey: "los_cristianos", EventType: "HIGH", Source: "open_meteo", EventTimeLocal: time.Date(2026, 3, 1, 12, 0, 0, 0, madrid)},
		},
	}

	collector := NewCollector(repo)
	collector.SetServingSource("hybrid")
	collector.SetFetchers(
		func(ctx context.Context, date time.Time, loc LocationConfig) ([]ExtremeEvent, string, error) {
			return []ExtremeEvent{
				{Type: "LOW", TimeLocal: time.Date(2026, 3, 1, 5, 57, 0, 0, madrid), HeightM: -1.2},
				{Type: "HIGH", TimeLocal: time.Date(2026, 3, 1, 12, 4, 0, 0, madrid), HeightM: 0.6},
			}, `{"source":"puertos"}`, nil
		},
		func(ctx context.Context, date time.Time, loc LocationConfig) ([]ExtremeEvent, string, error) {
			t.Fatal("open-meteo fetcher should not run when Puertos succeeds")
			return nil, "", nil
		},
	)

	if err := collector.CollectTides(context.Background(), "2026-03-01", "los_cristianos"); err != nil {
		t.Fatalf("CollectTides failed: %v", err)
	}
	if len(repo.events) != 2 {
		t.Fatalf("expected 2 replaced events, got %d", len(repo.events))
	}
	for _, e := range repo.events {
		if e.Source != "puertos" {
			t.Fatalf("expected puertos source after overwrite, got %s", e.Source)
		}
		if e.Confidence != 90 {
			t.Fatalf("expected confidence 90, got %d", e.Confidence)
		}
	}
}

func TestCollectTidesOpenMeteoOnlySkipsPuertos(t *testing.T) {
	repo := &fakeRepo{}
	collector := NewCollector(repo)
	collector.SetServingSource("open_meteo")

	madrid, _ := time.LoadLocation("Europe/Madrid")
	collector.SetFetchers(
		func(ctx context.Context, date time.Time, loc LocationConfig) ([]ExtremeEvent, string, error) {
			t.Fatal("puertos fetcher should not run in open_meteo mode")
			return nil, "", nil
		},
		func(ctx context.Context, date time.Time, loc LocationConfig) ([]ExtremeEvent, string, error) {
			return []ExtremeEvent{
				{Type: "LOW", TimeLocal: time.Date(2026, 3, 1, 6, 0, 0, 0, madrid), HeightM: -1.1},
				{Type: "HIGH", TimeLocal: time.Date(2026, 3, 1, 12, 0, 0, 0, madrid), HeightM: 0.5},
			}, `{"source":"open"}`, nil
		},
	)

	if err := collector.CollectTides(context.Background(), "2026-03-01", "los_cristianos"); err != nil {
		t.Fatalf("CollectTides failed: %v", err)
	}
	if len(repo.events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(repo.events))
	}
	for _, e := range repo.events {
		if e.Source != "open_meteo" {
			t.Fatalf("expected open_meteo source, got %s", e.Source)
		}
	}
}
