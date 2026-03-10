-- Safety migration for environments where another version 0003 was already applied.
-- Keeps creation idempotent.
CREATE TABLE IF NOT EXISTS pws_stations (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  station_id VARCHAR(64) NOT NULL,
  name VARCHAR(120) NOT NULL,
  lat DECIMAL(9,6) DEFAULT NULL,
  lon DECIMAL(9,6) DEFAULT NULL,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  display_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY pws_stations_station_id_unique (station_id),
  KEY pws_stations_active_order_index (is_active, display_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS pws_latest (
  station_ref_id BIGINT UNSIGNED NOT NULL,
  temp_c DECIMAL(5,1) DEFAULT NULL,
  humidity DECIMAL(5,1) DEFAULT NULL,
  obs_time_utc DATETIME DEFAULT NULL,
  fetched_at_utc DATETIME NOT NULL,
  stale TINYINT(1) NOT NULL DEFAULT 0,
  invalid TINYINT(1) NOT NULL DEFAULT 0,
  error_message VARCHAR(255) DEFAULT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (station_ref_id),
  KEY pws_latest_fetched_at_index (fetched_at_utc),
  CONSTRAINT pws_latest_station_ref_id_fk
    FOREIGN KEY (station_ref_id) REFERENCES pws_stations(id)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
