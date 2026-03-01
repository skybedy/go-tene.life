CREATE TABLE IF NOT EXISTS tide_events (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  date_local DATE NOT NULL COMMENT 'Local date (Europe/Madrid)',
  location_key VARCHAR(64) NOT NULL,
  event_type VARCHAR(8) NOT NULL COMMENT 'HIGH|LOW',
  event_time_local DATETIME NOT NULL COMMENT 'Local datetime (Europe/Madrid)',
  height_m DECIMAL(8,3) NOT NULL,
  source VARCHAR(32) NOT NULL COMMENT 'puertos|open_meteo',
  confidence TINYINT UNSIGNED NOT NULL,
  fetched_at DATETIME NOT NULL COMMENT 'UTC fetch timestamp',
  raw_json LONGTEXT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY tide_events_unique (date_local, location_key, event_type, event_time_local),
  KEY tide_events_lookup_idx (date_local, location_key, source, fetched_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
