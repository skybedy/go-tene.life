package tides

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/skybedy/laravel-tene.life/internal/models"
)

type Repository interface {
	GetTideEvents(ctx context.Context, dateLocal, locationKey string) ([]models.TideEvent, error)
	HasFreshPuertosData(ctx context.Context, dateLocal, locationKey string, sinceUTC time.Time) (bool, error)
	ReplaceTideEvents(ctx context.Context, dateLocal, locationKey string, events []models.TideEvent) error
}

type Collector struct {
	repo            Repository
	fetchPuertos    ExtremeFetcher
	fetchOpenMeteo  ExtremeFetcher
	freshnessWindow time.Duration
	servingSource   string
}

func NewCollector(repo Repository) *Collector {
	return &Collector{
		repo:            repo,
		fetchPuertos:    FetchPuertosExtremes,
		fetchOpenMeteo:  FetchOpenMeteoExtremes,
		freshnessWindow: 6 * time.Hour,
		servingSource:   ServingSource(),
	}
}

func (c *Collector) SetFetchers(puertos, openMeteo ExtremeFetcher) {
	if puertos != nil {
		c.fetchPuertos = puertos
	}
	if openMeteo != nil {
		c.fetchOpenMeteo = openMeteo
	}
}

func (c *Collector) SetServingSource(source string) {
	c.servingSource = normalizeServingSource(source)
}

func (c *Collector) CollectTides(ctx context.Context, dateLocal, locationKey string) error {
	locCfg, ok := ResolveLocation(locationKey)
	if !ok {
		return fmt.Errorf("unknown location key: %s", locationKey)
	}
	if strings.TrimSpace(dateLocal) == "" {
		tz, err := time.LoadLocation(locCfg.Timezone)
		if err != nil {
			return err
		}
		dateLocal = time.Now().In(tz).Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateLocal)
	if err != nil {
		return fmt.Errorf("invalid date %q: %w", dateLocal, err)
	}

	servingSource := normalizeServingSource(c.servingSource)
	freshnessSource := servingSource
	if servingSource == "hybrid" {
		freshnessSource = "puertos"
	}
	fresh, err := hasFreshSourceData(ctx, c.repo, dateLocal, locCfg.Key, freshnessSource, time.Now().UTC().Add(-c.freshnessWindow))
	if err != nil {
		return err
	}
	if fresh {
		return nil
	}

	if servingSource == "open_meteo" {
		openEvents, openRaw, openErr := c.fetchOpenMeteo(ctx, date, locCfg)
		if openErr != nil {
			log.Printf("tides source=open_meteo reason=fetch_failed date=%s location=%s error=%v", dateLocal, locCfg.Key, openErr)
			return openErr
		}

		if err := validateEventsForDate(dateLocal, locCfg.Timezone, openEvents); err != nil {
			log.Printf("tides source=open_meteo reason=validation_failed date=%s location=%s error=%v", dateLocal, locCfg.Key, err)
			return err
		}

		rows := toDBEvents(dateLocal, locCfg.Key, "open_meteo", 70, openRaw, openEvents)
		return c.repo.ReplaceTideEvents(ctx, dateLocal, locCfg.Key, rows)
	}

	puertosEvents, puertosRaw, err := c.fetchPuertos(ctx, date, locCfg)
	if err == nil {
		if err := validateEventsForDate(dateLocal, locCfg.Timezone, puertosEvents); err == nil {
			rows := toDBEvents(dateLocal, locCfg.Key, "puertos", 90, puertosRaw, puertosEvents)
			return c.repo.ReplaceTideEvents(ctx, dateLocal, locCfg.Key, rows)
		}
		log.Printf("tides fallback source=puertos reason=validation_failed date=%s location=%s error=%v", dateLocal, locCfg.Key, err)
	} else {
		log.Printf("tides fallback source=puertos reason=fetch_failed date=%s location=%s error=%v", dateLocal, locCfg.Key, err)
	}

	openEvents, openRaw, openErr := c.fetchOpenMeteo(ctx, date, locCfg)
	if openErr != nil {
		log.Printf("tides fallback source=open_meteo reason=fetch_failed date=%s location=%s error=%v", dateLocal, locCfg.Key, openErr)
		return openErr
	}

	if err := validateEventsForDate(dateLocal, locCfg.Timezone, openEvents); err != nil {
		log.Printf("tides fallback source=open_meteo reason=validation_failed date=%s location=%s error=%v", dateLocal, locCfg.Key, err)
		return err
	}

	rows := toDBEvents(dateLocal, locCfg.Key, "open_meteo", 70, openRaw, openEvents)
	return c.repo.ReplaceTideEvents(ctx, dateLocal, locCfg.Key, rows)
}

func ServingSource() string {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv("TIDES_SERVING_SOURCE")))
	return normalizeServingSource(raw)
}

func normalizeServingSource(source string) string {
	switch strings.TrimSpace(strings.ToLower(source)) {
	case "", "open_meteo", "open-meteo", "openmeteo":
		return "open_meteo"
	case "hybrid", "puertos_fallback", "puertos->open_meteo":
		return "hybrid"
	default:
		return "open_meteo"
	}
}

func hasFreshSourceData(ctx context.Context, repo Repository, dateLocal, locationKey, source string, sinceUTC time.Time) (bool, error) {
	if source == "puertos" {
		return repo.HasFreshPuertosData(ctx, dateLocal, locationKey, sinceUTC)
	}

	events, err := repo.GetTideEvents(ctx, dateLocal, locationKey)
	if err != nil {
		return false, err
	}
	for _, ev := range events {
		if ev.Source != source {
			continue
		}
		if ev.FetchedAt.UTC().After(sinceUTC) {
			return true, nil
		}
	}
	return false, nil
}

func toDBEvents(dateLocal, locationKey, source string, confidence int, raw string, events []ExtremeEvent) []models.TideEvent {
	fetchedAt := time.Now().UTC()
	out := make([]models.TideEvent, 0, len(events))
	for _, e := range events {
		out = append(out, models.TideEvent{
			DateLocal:      dateLocal,
			LocationKey:    locationKey,
			EventType:      strings.ToUpper(e.Type),
			EventTimeLocal: e.TimeLocal,
			HeightM:        e.HeightM,
			Source:         source,
			Confidence:     confidence,
			FetchedAt:      fetchedAt,
			RawJSON:        raw,
		})
	}
	return out
}

func validateEventsForDate(dateLocal, timezone string, events []ExtremeEvent) error {
	if len(events) < 2 {
		return fmt.Errorf("expected at least 2 events, got %d", len(events))
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	for _, e := range events {
		eventDate := e.TimeLocal.In(loc).Format("2006-01-02")
		if eventDate != dateLocal {
			return fmt.Errorf("event %s at %s is outside date %s", e.Type, e.TimeLocal.Format(time.RFC3339), dateLocal)
		}
		if e.Type != "HIGH" && e.Type != "LOW" {
			return fmt.Errorf("invalid event type %q", e.Type)
		}
	}
	return nil
}
