package store

import (
	"database/sql"
	"time"

	"github.com/skybedy/laravel-tene.life/internal/models"
)

func (s *WeatherStore) GetActivePWSStations() ([]models.PWSStation, error) {
	query := `SELECT id, station_id, name, lat, lon
		FROM pws_stations
		WHERE is_active = 1
		ORDER BY display_order ASC, name ASC`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stations := make([]models.PWSStation, 0)
	for rows.Next() {
		var id uint64
		var stationID, name string
		var lat, lon sql.NullFloat64
		if err := rows.Scan(&id, &stationID, &name, &lat, &lon); err != nil {
			return nil, err
		}

		station := models.PWSStation{
			ID:        id,
			StationID: stationID,
			Name:      name,
		}
		if lat.Valid {
			v := lat.Float64
			station.Lat = &v
		}
		if lon.Valid {
			v := lon.Float64
			station.Lon = &v
		}
		stations = append(stations, station)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stations, nil
}

func (s *WeatherStore) UpsertPWSLatest(rec models.PWSLatestRecord) error {
	query := `INSERT INTO pws_latest
		(station_ref_id, temp_c, humidity, obs_time_utc, fetched_at_utc, stale, invalid, error_message)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		temp_c = VALUES(temp_c),
		humidity = VALUES(humidity),
		obs_time_utc = VALUES(obs_time_utc),
		fetched_at_utc = VALUES(fetched_at_utc),
		stale = VALUES(stale),
		invalid = VALUES(invalid),
		error_message = VALUES(error_message)`

	_, err := s.DB.Exec(
		query,
		rec.StationRefID,
		nullableFloat(rec.TempC),
		nullableFloat(rec.Humidity),
		nullableTime(rec.ObsTimeUTC),
		rec.FetchedAtUTC.UTC(),
		rec.Stale,
		rec.Invalid,
		nullableString(rec.ErrorMessage),
	)
	return err
}

func (s *WeatherStore) GetPWSLatestPoints() ([]models.PWSMapPoint, error) {
	query := `SELECT
		s.station_id,
		s.name,
		s.lat,
		s.lon,
		l.temp_c,
		l.humidity,
		l.obs_time_utc,
		l.fetched_at_utc,
		l.stale,
		l.invalid,
		l.error_message
	FROM pws_stations s
	LEFT JOIN pws_latest l ON l.station_ref_id = s.id
	WHERE s.is_active = 1
	ORDER BY s.display_order ASC, s.name ASC`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := make([]models.PWSMapPoint, 0)
	for rows.Next() {
		var stationID, name string
		var lat, lon sql.NullFloat64
		var temp, humidity sql.NullFloat64
		var obsTime, fetchedAt sql.NullTime
		var stale, invalid sql.NullBool
		var errorMsg sql.NullString

		if err := rows.Scan(
			&stationID,
			&name,
			&lat,
			&lon,
			&temp,
			&humidity,
			&obsTime,
			&fetchedAt,
			&stale,
			&invalid,
			&errorMsg,
		); err != nil {
			return nil, err
		}

		point := models.PWSMapPoint{
			StationID: stationID,
			Name:      name,
			Stale:     true,
			Invalid:   false,
		}

		if lat.Valid {
			v := lat.Float64
			point.Lat = &v
		}
		if lon.Valid {
			v := lon.Float64
			point.Lon = &v
		}
		if temp.Valid {
			v := temp.Float64
			point.TempC = &v
		}
		if humidity.Valid {
			v := humidity.Float64
			point.Humidity = &v
		}
		if obsTime.Valid {
			point.ObsTimeUTC = obsTime.Time.UTC().Format(time.RFC3339)
		}
		if fetchedAt.Valid {
			point.FetchedAtUTC = fetchedAt.Time.UTC().Format(time.RFC3339)
		}
		if stale.Valid {
			point.Stale = stale.Bool
		}
		if invalid.Valid {
			point.Invalid = invalid.Bool
		}
		if errorMsg.Valid {
			point.Error = errorMsg.String
		}

		points = append(points, point)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return points, nil
}

func nullableFloat(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullableTime(v *time.Time) any {
	if v == nil {
		return nil
	}
	return v.UTC()
}

func nullableString(v string) any {
	if v == "" {
		return nil
	}
	return v
}
