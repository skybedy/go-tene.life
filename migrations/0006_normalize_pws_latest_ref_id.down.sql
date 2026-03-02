-- Rollback to legacy pws_latest shape keyed by station_id.

CREATE TABLE IF NOT EXISTS pws_latest_old (
  station_id VARCHAR(64) NOT NULL,
  temp_c DECIMAL(5,1) DEFAULT NULL,
  humidity DECIMAL(5,1) DEFAULT NULL,
  lat DECIMAL(9,6) DEFAULT NULL,
  lon DECIMAL(9,6) DEFAULT NULL,
  obs_time_utc DATETIME DEFAULT NULL,
  fetched_at_utc DATETIME NOT NULL,
  stale TINYINT(1) NOT NULL DEFAULT 0,
  invalid TINYINT(1) NOT NULL DEFAULT 0,
  error_message VARCHAR(255) DEFAULT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (station_id),
  KEY pws_latest_fetched_at_index (fetched_at_utc),
  CONSTRAINT pws_latest_station_id_fk
    FOREIGN KEY (station_id) REFERENCES pws_stations(station_id)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO pws_latest_old (
  station_id, temp_c, humidity, lat, lon, obs_time_utc, fetched_at_utc, stale, invalid, error_message, updated_at
)
SELECT
  s.station_id,
  l.temp_c,
  l.humidity,
  s.lat,
  s.lon,
  l.obs_time_utc,
  l.fetched_at_utc,
  l.stale,
  l.invalid,
  l.error_message,
  l.updated_at
FROM pws_latest l
JOIN pws_stations s ON s.id = l.station_ref_id
ON DUPLICATE KEY UPDATE
  temp_c = VALUES(temp_c),
  humidity = VALUES(humidity),
  lat = VALUES(lat),
  lon = VALUES(lon),
  obs_time_utc = VALUES(obs_time_utc),
  fetched_at_utc = VALUES(fetched_at_utc),
  stale = VALUES(stale),
  invalid = VALUES(invalid),
  error_message = VALUES(error_message),
  updated_at = VALUES(updated_at);

DROP TABLE pws_latest;
RENAME TABLE pws_latest_old TO pws_latest;
