-- 1) Ověření, že legacy sloupec už neexistuje
SELECT COUNT(*) AS legacy_column_exists
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'weather_daily'
  AND COLUMN_NAME = 'sea_temperature';

-- 2) Kolik řádků je v nové tabulce celkem
SELECT COUNT(*) AS water_temperatures_total
FROM water_temperatures;

-- 3) Kolik backfill řádků bylo vloženo
SELECT COUNT(*) AS backfilled_rows
FROM water_temperatures
WHERE source = 'legacy_backfill';

-- 4) Které legacy řádky (podle reference) chybí po backfillu
SELECT wd.id, wd.date
FROM weather_daily wd
LEFT JOIN water_temperatures wt ON wt.legacy_weather_daily_id = wd.id
WHERE wt.id IS NULL
ORDER BY wd.date ASC;

-- 5) Duplicity (neměly by existovat)
SELECT legacy_weather_daily_id, COUNT(*) AS cnt
FROM water_temperatures
WHERE legacy_weather_daily_id IS NOT NULL
GROUP BY legacy_weather_daily_id
HAVING COUNT(*) > 1;

-- 6) Nejnovější měření teploty vody
SELECT id, measured_at, temperature, source, note, created_at, updated_at
FROM water_temperatures
ORDER BY measured_at DESC
LIMIT 1;

-- 7) Rychlá kontrola více měření za den
SELECT DATE(measured_at) AS measured_day, COUNT(*) AS measurements_in_day
FROM water_temperatures
GROUP BY DATE(measured_at)
HAVING COUNT(*) > 1
ORDER BY measured_day DESC;
