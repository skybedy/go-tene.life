# go-tene.life

## Database

This project uses versioned SQL migrations for MariaDB.

- Schema snapshot: `db/schema.sql`
- Migration files: `migrations/*.up.sql` and `migrations/*.down.sql`
- Migrator: `golang-migrate`

### Prerequisites

- MariaDB running and reachable with credentials from `.env`
- Go installed (used to run the migrate CLI)
- `mysqldump` installed (for schema export)

### Run migrations

```bash
make migrate-up
```

Rollback last migration:

```bash
make migrate-down
```

Show current migration version:

```bash
make migrate-status
```

### Generate schema snapshot

```bash
make dump-schema
```

This exports structure-only SQL for app tables (`weather*`, `tide_events`, `pws_*`) into `db/schema.sql`.

### Recommended workflow for DB changes

1. Add a new migration pair in `migrations/` with the next version number.
2. Apply it locally using `make migrate-up`.
3. Refresh snapshot using `make dump-schema`.
4. Commit migration files and `db/schema.sql` together.

## Waves Collector

Measured wave data is collected from Puertos del Estado (PORTUS), station `2446` (Tenerife Sur), and cached to JSON.

- In `APP_ENV=local` / `dev` / `development`, startup collectors are disabled by default to avoid hitting external APIs on each local restart.
- To force them on locally, set `ENABLE_EXTERNAL_COLLECTORS=1`.
- Collector runs automatically inside the app (immediately on startup, then periodically).
- Default interval: every 15 minutes
- Optional env override: `WAVES_COLLECT_INTERVAL_MINUTES=15`
- Manual command: `tenelife collect:waves`
- Output cache file: `data/waves_latest.json`
- Source is fetched only by collector; request handlers read only the JSON cache.

Run manually:

```bash
go run . collect:waves
```

or with built binary:

```bash
./tenelife collect:waves
```

Optional cron fallback (only if you do not want the in-app collector):

```cron
*/15 * * * * cd /path/to/go-tene.life && ./tenelife collect:waves >> /var/log/tenelife-waves.log 2>&1
```

## Water Quality Collector

Official bathing water quality is collected from IDECanarias / GRAFCAN WMS (`ZB_PM`, layer `PM`) and cached to JSON.

- Manual command: `tenelife collect:water`
- Output cache file: `data/water_quality_latest.json`
- API uses cache only: `/api/home` reads this file and never calls external source during request handling.

Run manually:

```bash
go run . collect:water
```

or with built binary:

```bash
./tenelife collect:water
```

Cron example (once per day):

```cron
15 6 * * * cd /path/to/go-tene.life && ./tenelife collect:water >> /var/log/tenelife-water.log 2>&1
```

In-app scheduler (recommended, no Linux cron/systemd needed):

- `WATER_COLLECT_INTERVAL_MINUTES=1440` (default: once per 24h)
- collector runs automatically on app start and then in this interval

## Tenerife PWS Temperature Map

Current temperatures on Tenerife are loaded from The Weather Company PWS API and stored in DB cache tables.

- API key env: `WEATHER_COM_API_KEY`
- Collector interval env: `PWS_COLLECT_INTERVAL_MINUTES=10`
- WU cache env: `WU_CACHE_TTL_SECONDS=60`
- WU rate-limit env: `WU_RATELIMIT_PER_MIN=25`, `WU_RATELIMIT_BURST=5`
- WU stale-fallback max age: `WU_STALE_FALLBACK_MAX_AGE_SECONDS=120`
- Manual command: `go run . collect:pws`
- API endpoint: `/api/tenerife/pws-latest`
- Debug usage endpoint: `/debug/wu-usage`
- Page: `/tenerife/teploty` (also locale-prefixed variants)

### DB tables

- `pws_stations`: station configuration (`station_id`, `name`, optional `lat`/`lon`, `is_active`, `display_order`)
- `pws_latest`: latest fetched values per station (FK `station_ref_id -> pws_stations.id`, values: `temp_c`, `humidity`, `obs_time_utc`, `fetched_at_utc`, `stale`, `invalid`)

Example station inserts:

```sql
INSERT INTO pws_stations (station_id, name, lat, lon, is_active, display_order) VALUES
('ICANARIA12', 'Los Cristianos', 28.0436, -16.7215, 1, 10),
('ICANARIA45', 'Costa Adeje', 28.0900, -16.7350, 1, 20);
```

Cron example (every 10 minutes):

```cron
*/10 * * * * cd /path/to/go-tene.life && ./tenelife collect:pws >> /var/log/tenelife-pws.log 2>&1
```

## Tides API

Daily tide extremes (high/low with time and height) are stored in DB table `tide_events`.

- Endpoint: `/api/tides?date=YYYY-MM-DD&loc=los_cristianos`
- Default serving source: `open_meteo` (`TIDES_SERVING_SOURCE=open_meteo`)
- Optional hybrid mode: `TIDES_SERVING_SOURCE=hybrid` (Puertos first, fallback Open-Meteo)
- If data is missing, endpoint triggers synchronous collect with short timeout and may return `503` (`try_later`).

## Copernicus Sea Temperature Viewers

Static pages for pre-generated Copernicus outputs:

- Tenerife: `/atlantic-sea-temperature-tenerife`
- Canary Islands: `/atlantic-sea-temperature-canary-islands`

Legacy aliases still work via redirect:

- `/teplota-more` -> `/atlantic-sea-temperature-tenerife`
- `/teplota-more-kanary` -> `/atlantic-sea-temperature-canary-islands`

Viewer loads monthly manifest by selected period from:

- `/data/copernicus/sea-temp/tenerife/YYYY/MM/manifest.json`
- `/data/copernicus/sea-temp/canary/YYYY/MM/manifest.json`

Required structure in `public/`:

```text
public/data/copernicus/sea-temp/{tenerife|canary}/YYYY/MM/
  manifest.json
  YYYY-MM-DD.png (or .jpg/.svg)
  ...
```

### How to copy exports from Meteodata (manual phase-1 sync)

Copy generated month export from Meteodata into this repo under `public/data/...`.

Example:

```bash
# from go-tene.life root
mkdir -p public/data/copernicus/sea-temp/tenerife/2026/04
cp -R /path/to/meteodata/exports/copernicus/sea-temp/tenerife/2026/04/* public/data/copernicus/sea-temp/tenerife/2026/04/

mkdir -p public/data/copernicus/sea-temp/canary/2026/04
cp -R /path/to/meteodata/exports/copernicus/sea-temp/canary/2026/04/* public/data/copernicus/sea-temp/canary/2026/04/
```

After copy, open one of the viewer URLs above; page loads `manifest.json` via `fetch()` and renders slider/playback from static files only.


> Note: Real Copernicus exports from Meteodata (especially PNG/JPG/SVG/MP4) should **not** be committed in PRs.
> Copy them to `public/data/...` only for local/server static deployment. Keep repo demo assets text-based (e.g. SVG + manifest).
