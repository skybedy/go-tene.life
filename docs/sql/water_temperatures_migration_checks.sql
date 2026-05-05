-- 1) Kolik historických hodnot teploty vody je v legacy weather_daily
SELECT COUNT(*) AS legacy_non_null_count
FROM weather_daily
WHERE sea_temperature IS NOT NULL;

-- 2) Kolik řádků je v nové tabulce celkem
SELECT COUNT(*) AS water_temperatures_total
FROM water_temperatures;

-- 3) Kolik backfill řádků bylo vloženo
SELECT COUNT(*) AS backfilled_rows
FROM water_temperatures
WHERE source = 'legacy_backfill';

-- 4) Které legacy hodnoty chybí po backfillu
SELECT wd.id, wd.date, wd.sea_temperature
FROM weather_daily wd
LEFT JOIN water_temperatures wt ON wt.legacy_weather_daily_id = wd.id
WHERE wd.sea_temperature IS NOT NULL
  AND wt.id IS NULL
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
