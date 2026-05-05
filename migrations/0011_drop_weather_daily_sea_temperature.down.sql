ALTER TABLE weather_daily
  ADD COLUMN sea_temperature DECIMAL(5,1) NULL COMMENT 'Sea water temperature in °C (manually measured)' AFTER date;
