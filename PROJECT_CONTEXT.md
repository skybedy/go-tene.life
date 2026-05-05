# PROJECT_CONTEXT.md

## Projekt

- Název: `go-tene.life`
- Účel: webová aplikace pro informace související s Tenerife (počasí, vlny, kvalita vody, příliv/odliv, PWS teploty).

## Aktuální stav

- Typ: rozpracovaný existující projekt.
- Git: aktuální větev `feature/water-temperature-refactor`; pracovní strom obsahuje lokální úpravu `views/index.html` (UI swap barev ve weather boxu, zatím bez commitu).
- AI kontextové soubory byly vytvořeny v tomto kroku.

## Technologie

- Backend: Go `1.24`, Echo v4.
- Databáze: MariaDB/MySQL přes `github.com/go-sql-driver/mysql`.
- Konfigurace: `.env` přes `github.com/joho/godotenv`.
- Frontend: server-side HTML šablony (`views/`) + Vanilla JS (`public/js`).
- Styling: Tailwind CSS 4 (`@tailwindcss/cli`), zdroj `resources/css/app.css`, výstup `public/css/app.css`.
- Migrace: SQL migrace v `migrations/` přes `golang-migrate` (Makefile).
- Sea temperature redesign (in progress): nové měření jde do `water_temperatures` (timestamp + hodnota), legacy `weather_daily.sea_temperature` je zatím ponechán kvůli bezpečné etapizaci.
- Backfill pravidlo: historické hodnoty z `weather_daily.sea_temperature` se mapují na čas `10:00:00` (UTC) jako default pro legacy denní záznamy.

## Důležité adresáře a soubory

- Vstup aplikace: `main.go`
- Interní logika: `internal/` (např. `web`, `store`, `waves`, `water`, `pws`, `tides`)
- Šablony: `views/`
- Statické soubory: `public/`
- CSS zdroj: `resources/css/app.css`
- DB migrace: `migrations/`
- DB snapshot: `db/schema.sql`
- SQL ověření migrace: `docs/sql/water_temperatures_migration_checks.sql`
- Build/ops: `Makefile`, `deploy.sh`

## Spuštění

- Aplikace: `go run .`
- Kolektory ručně:
  - `go run . collect:waves`
  - `go run . collect:water`
  - `go run . collect:pws`

## Build a testy

- Go testy: `go test ./...`
- Go build: `go build ./...`
- CSS build: `npm run build`
- CSS watch: `npm run dev`
- DB migrace: `make migrate-up`, `make migrate-down`, `make migrate-status`
- Export DB schématu: `make dump-schema`

## Poznámky pro další práci

- `node_modules/` je přítomné v projektu.
- V `data/` jsou cache JSON soubory používané kolektory.
- Konkrétní deployment workflow je zatím nezjištěno (k doplnění).
- V homepage weather boxu proběhl swap barev (to, co bylo `text-white/90`, je nyní `text-orange-300`; původní `text-orange-300` je nyní `text-white/90`).
