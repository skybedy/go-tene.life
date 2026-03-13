CREATE TABLE IF NOT EXISTS weather (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  measured_at DATETIME NOT NULL COMMENT 'Time when the weather data was measured',
  temperature DECIMAL(5,1) NOT NULL COMMENT 'Temperature in °C',
  pressure DECIMAL(7,1) NOT NULL COMMENT 'Pressure in hPa',
  humidity DECIMAL(5,1) NOT NULL COMMENT 'Humidity in %',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weather_daily (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  date DATE NOT NULL COMMENT 'Date of the measurement',
  sea_temperature DECIMAL(5,1) NULL COMMENT 'Sea water temperature in °C (manually measured)',
  avg_temperature DECIMAL(5,1) NULL,
  min_temperature DECIMAL(5,1) NULL,
  max_temperature DECIMAL(5,1) NULL,
  avg_pressure DECIMAL(7,1) NULL,
  min_pressure DECIMAL(7,1) NULL,
  max_pressure DECIMAL(7,1) NULL,
  avg_humidity DECIMAL(5,1) NULL,
  min_humidity DECIMAL(5,1) NULL,
  max_humidity DECIMAL(5,1) NULL,
  samples_count INT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY weather_daily_date_unique (date),
  KEY weather_daily_date_index (date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weather_hourly (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  date DATE NOT NULL COMMENT 'Date of the measurement',
  hour TINYINT NOT NULL COMMENT 'Hour (0-23)',
  avg_temperature DECIMAL(5,1) NOT NULL COMMENT 'Average temperature in °C',
  avg_pressure DECIMAL(7,1) NOT NULL COMMENT 'Average pressure in hPa',
  avg_humidity DECIMAL(5,1) NOT NULL COMMENT 'Average humidity in %',
  samples_count INT NOT NULL COMMENT 'Number of measurements used for average',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY weather_hourly_date_hour_unique (date, hour),
  KEY weather_hourly_date_index (date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weather_monthly (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  year INT NOT NULL COMMENT 'Year (e.g., 2025)',
  month TINYINT NOT NULL COMMENT 'Month (1-12)',
  avg_sea_temperature DECIMAL(5,1) NULL COMMENT 'Average sea temperature in °C',
  avg_temperature DECIMAL(5,1) NOT NULL COMMENT 'Average temperature in °C',
  min_temperature DECIMAL(5,1) NOT NULL COMMENT 'Minimum temperature in °C',
  max_temperature DECIMAL(5,1) NOT NULL COMMENT 'Maximum temperature in °C',
  avg_pressure DECIMAL(7,1) NOT NULL COMMENT 'Average pressure in hPa',
  min_pressure DECIMAL(7,1) NOT NULL COMMENT 'Minimum pressure in hPa',
  max_pressure DECIMAL(7,1) NOT NULL COMMENT 'Maximum pressure in hPa',
  avg_humidity DECIMAL(5,1) NOT NULL COMMENT 'Average humidity in %',
  min_humidity DECIMAL(5,1) NOT NULL COMMENT 'Minimum humidity in %',
  max_humidity DECIMAL(5,1) NOT NULL COMMENT 'Maximum humidity in %',
  samples_count INT NOT NULL COMMENT 'Number of measurements used',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY weather_monthly_year_month_unique (year, month),
  KEY weather_monthly_year_month_index (year, month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weather_weekly (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  year INT NOT NULL COMMENT 'Year (e.g., 2025)',
  week TINYINT NOT NULL COMMENT 'ISO week number (1-53)',
  week_start DATE NOT NULL COMMENT 'Monday of the week',
  week_end DATE NOT NULL COMMENT 'Sunday of the week',
  avg_sea_temperature DECIMAL(5,1) NULL COMMENT 'Average sea temperature in °C',
  avg_temperature DECIMAL(5,1) NOT NULL COMMENT 'Average temperature in °C',
  min_temperature DECIMAL(5,1) NOT NULL COMMENT 'Minimum temperature in °C',
  max_temperature DECIMAL(5,1) NOT NULL COMMENT 'Maximum temperature in °C',
  avg_pressure DECIMAL(7,1) NOT NULL COMMENT 'Average pressure in hPa',
  min_pressure DECIMAL(7,1) NOT NULL COMMENT 'Minimum pressure in hPa',
  max_pressure DECIMAL(7,1) NOT NULL COMMENT 'Maximum pressure in hPa',
  avg_humidity DECIMAL(5,1) NOT NULL COMMENT 'Average humidity in %',
  min_humidity DECIMAL(5,1) NOT NULL COMMENT 'Minimum humidity in %',
  max_humidity DECIMAL(5,1) NOT NULL COMMENT 'Maximum humidity in %',
  samples_count INT NOT NULL COMMENT 'Number of measurements used',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY weather_weekly_year_week_unique (year, week),
  KEY weather_weekly_year_week_index (year, week)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
