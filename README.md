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

This exports structure-only SQL for app tables (`weather*`, `tide_events`) into `db/schema.sql`.

### Recommended workflow for DB changes

1. Add a new migration pair in `migrations/` with the next version number.
2. Apply it locally using `make migrate-up`.
3. Refresh snapshot using `make dump-schema`.
4. Commit migration files and `db/schema.sql` together.

## Waves Collector

Measured wave data is collected from Puertos del Estado (PORTUS), station `2446` (Tenerife Sur), and cached to JSON.

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

## Tides API

Daily tide extremes (high/low with time and height) are stored in DB table `tide_events`.

- Endpoint: `/api/tides?date=YYYY-MM-DD&loc=los_cristianos`
- Default serving source: `open_meteo` (`TIDES_SERVING_SOURCE=open_meteo`)
- Optional hybrid mode: `TIDES_SERVING_SOURCE=hybrid` (Puertos first, fallback Open-Meteo)
- If data is missing, endpoint triggers synchronous collect with short timeout and may return `503` (`try_later`).
