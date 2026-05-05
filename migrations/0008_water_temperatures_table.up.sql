-- Creates dedicated point-in-time sea water temperature storage.
-- Keeps legacy weather_daily.sea_temperature untouched for phased rollout.
CREATE TABLE IF NOT EXISTS water_temperatures (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  measured_at DATETIME NOT NULL COMMENT 'Exact sea temperature measurement timestamp (UTC)',
  temperature DECIMAL(5,1) NOT NULL COMMENT 'Sea water temperature in °C',
  source VARCHAR(32) NOT NULL DEFAULT 'manual' COMMENT 'Measurement source (e.g. manual)',
  note VARCHAR(255) NULL,
  legacy_weather_daily_id BIGINT UNSIGNED NULL COMMENT 'Optional link to migrated weather_daily row',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY water_temperatures_legacy_daily_unique (legacy_weather_daily_id),
  KEY water_temperatures_measured_at_idx (measured_at),
  KEY water_temperatures_source_measured_at_idx (source, measured_at),
  CONSTRAINT water_temperatures_legacy_daily_fk
    FOREIGN KEY (legacy_weather_daily_id) REFERENCES weather_daily(id)
    ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
