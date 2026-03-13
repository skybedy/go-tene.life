ALTER TABLE weather_monthly
  ADD COLUMN avg_sea_temperature DECIMAL(5,1) NULL COMMENT 'Average sea temperature in °C' AFTER month;

ALTER TABLE weather_weekly
  ADD COLUMN avg_sea_temperature DECIMAL(5,1) NULL COMMENT 'Average sea temperature in °C' AFTER week_end;
