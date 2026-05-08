# PROJECT_CONTEXT.md

## Projekt

- Název: `go-tene.life`
- Účel: webová aplikace pro informace související s Tenerife (počasí, vlny, kvalita vody, příliv/odliv, PWS teploty).

## Aktuální stav

- Typ: rozpracovaný existující projekt.
- Git: aktuální větev `feature/sound-segments-lesson-02`.
- Všechna slovíčka (lekce 1–26 i velké soubory 1–250, 251–500) mají vygenerované časované segmenty v `public/sounds/index.json`.

## Technologie

- Backend: Go `1.24`, Echo v4.
- Databáze: MariaDB/MySQL přes `github.com/go-sql-driver/mysql`.
- Konfigurace: `.env` přes `github.com/joho/godotenv`.
- Frontend: server-side HTML šablony (`views/`) + Vanilla JS (`public/js`).
- Styling: Tailwind CSS 4 (`@tailwindcss/cli`), zdroj `resources/css/app.css`, výstup `public/css/app.css`.
- Migrace: SQL migrace v `migrations/` přes `golang-migrate` (Makefile).
- Sea temperature redesign (active): nové měření jde do `water_temperatures` (timestamp + hodnota), legacy sloupec `weather_daily.sea_temperature` je určen k odstranění migrací `0011`.
- Backfill pravidlo: historické hodnoty z `weather_daily.sea_temperature` se mapují na čas `10:00:00` (UTC) jako default pro legacy denní záznamy.
- Denní statistiky (`/statistics/daily`) už berou teplotu moře z `water_temperatures` (nejnovější měření v rámci dne) a v UI zobrazují i timestamp poslední aktualizace.
- Na `/statistics/daily` je doplněný fallback: pokud pro dnešní datum chybí teplota moře, API vrátí poslední dostupné měření; UI pak zobrazuje `Dnes HH:MM` jen pokud měření opravdu pochází z dnešního dne, jinak zobrazí celé datum+čas uložené v `water_temperatures`.
- Na `/statistics/monthly` a dalších agregačních stránkách byl opraven sea graf: už nebere syrovou historii z `/api/water-temperatures/history`, ale agregační dataset `sea_temperature` z `/api/weather/{weekly|monthly|annual}`.
- Na `/statistics/recent` se ve sloupci moře zobrazuje i čas měření (`SeaMeasuredAt`) ve formátu `H:MM` bez počáteční nuly.
- Na `/spanelsko-ceska-slovicka` je pro volbu `Přehrát všechno za sebou` nově připravený dedikovaný track `spanelsko_ceska_slovicka_1_500.mp3` s 500 časovanými segmenty, takže se během přehrávání synchronně zobrazuje španělský i český text.

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
