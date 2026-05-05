-- Backfills historical sea temperature values from weather_daily into water_temperatures.
-- Historical fallback time is set to 10:00:00 UTC for daily-only legacy records.
-- Insert is idempotent thanks to unique constraint on legacy_weather_daily_id.
INSERT INTO water_temperatures (measured_at, temperature, source, note, legacy_weather_daily_id)
SELECT
  STR_TO_DATE(CONCAT(wd.date, ' 10:00:00'), '%Y-%m-%d %H:%i:%s') AS measured_at,
  wd.sea_temperature,
  'legacy_backfill' AS source,
  'Backfilled from weather_daily.sea_temperature (default 10:00:00 UTC)' AS note,
  wd.id AS legacy_weather_daily_id
FROM weather_daily wd
WHERE wd.sea_temperature IS NOT NULL
ON DUPLICATE KEY UPDATE legacy_weather_daily_id = legacy_weather_daily_id;
