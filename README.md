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

This exports structure-only SQL for app tables (`weather*`) into `db/schema.sql`.

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

## Tenerife PWS Temperature Map

Current temperatures on Tenerife are loaded from The Weather Company PWS API and stored in DB cache tables.

- API key env: `WEATHER_COM_API_KEY`
- Collector interval env: `PWS_COLLECT_INTERVAL_MINUTES=10`
- Manual command: `go run . collect:pws`
- API endpoint: `/api/tenerife/pws-latest`
- Page: `/tenerife/teploty` (also locale-prefixed variants)

### DB tables

- `pws_stations`: station configuration (`station_id`, `name`, optional `lat`/`lon`, `is_active`, `display_order`)
- `pws_latest`: latest fetched values per station (`temp_c`, `humidity`, `obs_time_utc`, `fetched_at_utc`, `stale`, `invalid`)

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
