package pws

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type countingRoundTripper struct {
	calls int64
	body  string
}

func (rt *countingRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	atomic.AddInt64(&rt.calls, 1)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(rt.body)),
		Header:     make(http.Header),
	}, nil
}

func TestWUClientCacheHitAvoidsSecondOutboundCall(t *testing.T) {
	t.Setenv("WEATHER_COM_API_KEY", "test-key")

	rt := &countingRoundTripper{
		body: `{"observations":[{"stationID":"ISTATION1","obsTimeUtc":"2026-03-02T15:00:00Z","epoch":1772463600,"lat":28.05,"lon":-16.71,"humidity":60,"metric":{"temp":22.1}}]}`,
	}
	client := newWUClient(
		wuClientConfig{cacheTTL: 60 * time.Second, rateLimitPerMin: 100, rateLimitBurst: 10},
		&http.Client{Transport: rt, Timeout: 2 * time.Second},
	)

	_, meta1, err := client.fetchCurrent(context.Background(), "ISTATION1", "m", "json")
	if err != nil {
		t.Fatalf("first fetch failed: %v", err)
	}
	if meta1.CacheHit {
		t.Fatalf("first fetch should not be cache hit")
	}

	_, meta2, err := client.fetchCurrent(context.Background(), "ISTATION1", "m", "json")
	if err != nil {
		t.Fatalf("second fetch failed: %v", err)
	}
	if !meta2.CacheHit {
		t.Fatalf("second fetch should be cache hit")
	}

	if got := atomic.LoadInt64(&rt.calls); got != 1 {
		t.Fatalf("expected 1 outbound call, got %d", got)
	}
}

func TestWUClientRateLimitReturns429WithoutCacheAndUsesCacheFallback(t *testing.T) {
	t.Setenv("WEATHER_COM_API_KEY", "test-key")

	rt := &countingRoundTripper{
		body: `{"observations":[{"stationID":"ISTATION1","obsTimeUtc":"2026-03-02T15:00:00Z","epoch":1772463600,"lat":28.05,"lon":-16.71,"humidity":60,"metric":{"temp":22.1}}]}`,
	}
	client := newWUClient(
		wuClientConfig{cacheTTL: 1 * time.Millisecond, rateLimitPerMin: 1, rateLimitBurst: 1},
		&http.Client{Transport: rt, Timeout: 2 * time.Second},
	)

	// Consume the single token and cache station A response.
	if _, _, err := client.fetchCurrent(context.Background(), "ISTATION_A", "m", "json"); err != nil {
		t.Fatalf("first fetch failed: %v", err)
	}

	// Different station has no cache and should be rate-limited.
	if _, _, err := client.fetchCurrent(context.Background(), "ISTATION_B", "m", "json"); err == nil {
		t.Fatalf("expected rate-limited error for station without cache")
	}

	time.Sleep(2 * time.Millisecond) // expire fresh-cache window for station A

	// Same station should use stale cached fallback when rate-limited.
	_, meta, err := client.fetchCurrent(context.Background(), "ISTATION_A", "m", "json")
	if err != nil {
		t.Fatalf("expected cache fallback for rate-limited station A, got err: %v", err)
	}
	if !meta.CacheHit {
		t.Fatalf("expected cache fallback to be marked as cache_hit")
	}

	if got := atomic.LoadInt64(&rt.calls); got != 1 {
		t.Fatalf("expected only 1 outbound call total, got %d", got)
	}
}
