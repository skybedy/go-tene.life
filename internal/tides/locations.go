package tides

import (
	"os"
	"strconv"
	"strings"
)

type LocationConfig struct {
	Key              string
	Lat              float64
	Lon              float64
	Timezone         string
	PuertosStationID int
	PuertosParam     string
}

var locations = map[string]LocationConfig{
	"los_cristianos": {
		Key:              "los_cristianos",
		Lat:              28.0436,
		Lon:              -16.7215,
		Timezone:         "Atlantic/Canary",
		PuertosStationID: 2446,
		PuertosParam:     "Ha",
	},
}

func ResolveLocation(key string) (LocationConfig, bool) {
	k := strings.TrimSpace(strings.ToLower(key))
	if k == "" {
		k = "los_cristianos"
	}
	loc, ok := locations[k]
	if !ok {
		return LocationConfig{}, false
	}

	if loc.Key == "los_cristianos" {
		if v, ok := envFloat("TIDE_LAT"); ok {
			loc.Lat = v
		}
		if v, ok := envFloat("TIDE_LON"); ok {
			loc.Lon = v
		}
		if tz := strings.TrimSpace(os.Getenv("TIDE_TIMEZONE")); tz != "" {
			loc.Timezone = tz
		}
	}

	return loc, ok
}

func envFloat(name string) (float64, bool) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
