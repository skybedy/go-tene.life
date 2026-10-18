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

This exports structure-only SQL for app tables (`weather*`, `tide_events`, `pws_*`, `news_*`) into `db/schema.sql`.

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

## Zprávy z Tenerife

Sekce je česky na `/zpravy`, detail na `/zpravy/{slug}`. Filtr používá
`?kategorie=doprava`, stránkování `?page=2` (12 článků na stránku). Homepage
zobrazuje poslední 3 publikované zprávy. `/sitemap.xml` zahrnuje publikované
články a stávající hlavní stránky, `/robots.txt` na ni odkazuje.

### Nasazení a první článek

1. Zkontrolujte `make migrate-status`. Nová migrace `0012_news.up.sql` vytváří
   pouze `news_categories` a `news_articles` a vkládá pět výchozích kategorií.
   Pokud jste na verzi 11, spusťte `make migrate-up`. U starší databáze nejprve
   prověřte předchozí migrace (zejména existující destruktivní migraci 0011).
   Při instalaci od nuly je navíc potřeba samostatně vyřešit existující
   migraci 0006 (kolize FK po 0005 a očekávání legacy `station_id`) a 0007
   (duplicitní sloupce již vytvořené v 0001).
   Aplikace sama migrace nespouští. Rollback 0012 odstraňuje články i kategorie.
2. Vytvořte bcrypt hash hesla pomocí pomocného příkazu (heslo se nezapisuje
   do historie příkazů ani argumentů procesu):

   ```bash
   read -rs -p 'Heslo redaktora (6–72 bajtů): ' editorpassword
   printf '\n'
   printf '%s' "$editorpassword" | go run ./cmd/news-password
   unset editorpassword
   ```

3. Do `.env` doplňte `SITE_URL=https://tene.life`, `NEWS_ADMIN_USERNAME` a
   `NEWS_ADMIN_PASSWORD_HASH='vygenerovaný bcrypt hash'`. Hash **musí být
   v jednoduchých uvozovkách**, aby `$` neinterpretoval shell ani dotenv.
   Prázdné obě přihlašovací hodnoty administraci vypínají (404).
   Pro lokální HTTP použijte `SITE_URL=http://localhost:8080`.
4. Spusťte `npm run build`, potom `go build -o tenelife .` a nasaďte/restartujte
   aplikaci stávajícím postupem. CSS a šablony jsou součástí Go binárky.
5. Otevřete `/admin/zpravy` a přihlaste se v dialogu prohlížeče. Klikněte na
   **Vytvořit článek**, vyplňte titulek, kategorii, perex a obsah. Slug se
   navrhne automaticky. Obsah používá Markdown; lišta editoru vkládá
   mezititulky, tučné písmo, kurzívu, seznamy, citace a odkazy. Raw HTML a
   nebezpečná URL se nevykreslí. Tlačítko náhledu používá stejný renderer jako
   veřejný detail článku.
6. Volitelně vyplňte původní zdroj (název i URL), obrázek a SEO titulek/popisek.
   Obrázek musí být dostupný přes HTTPS nebo `/images/…`; jeho soubor nahrajte
   stávajícím serverovým postupem do `public/images/`. Upload editor zatím nemá.
   Přidejte také popis obrázku. Prázdné SEO hodnoty převezmou titulek a perex.
7. Uložte koncept nebo použijte **Zobrazit náhled bez uložení**. Náhled
   zahrnuje aktuální obsah formuláře a nic neukládá. Pro publikaci vyberte
   **Publikovat na webu** a **Uložit článek**.

Datum první publikace se nastavuje automaticky v UTC, zobrazuje v čase
Tenerife a zůstává při opětovném publikování stejné. Od první publikace je slug
neměnný. Stažení z webu: nastavte **Koncept / stáhnout z webu** a uložte.
Smazání má samostatnou potvrzovací stránku. Souběžné úpravy se kontrolují
pomocí revize, aby starší formulář nepřepsal novější obsah.

### Provoz a zabezpečení

- Jediný redakční účet používá HTTP Basic autentizaci s bcrypt hashem.
  Nginx musí v produkci vynucovat HTTPS, správně předávat `Host` a
  `X-Forwarded-Proto` a mít backend dostupný pouze za proxy. Přihlášení
  prohlížeč uchovává; pro redakci na sdíleném počítači používejte soukromé okno
  a po skončení ho zavřete. Změna hashe a restart zneplatní staré heslo.
- Administrace má CSRF ochranu Echo (kontrola Fetch Metadata, u starších
  klientů token), cookie HttpOnly/SameSite=Strict/Secure na HTTPS, limit těla
  požadavku a limit ověřování přihlášení (20 požadavků naráz, obnova 1 za 2 s).
- SQL používá parametry, HTML šablony automaticky escapují text. Zdrojové URL
  jsou omezené na HTTP(S), obrázky na HTTPS nebo lokální `/images/…`.
- Koncepty, stažené a budoucí články se nevracejí veřejnými dotazy. Administrace
  i náhled mají `no-store` a `noindex`; veřejné zprávy vyžadují revalidaci.
  Pokud máte vlastní Nginx/CDN cache, respektujte tyto hlavičky a necachujte
  `/admin/`; nepoužívejte pro zprávy pravidlo, které ignoruje Cache-Control.
- Nginx musí předávat `/zpravy`, `/admin/zpravy`, `/sitemap.xml` a `/robots.txt`
  aplikaci; případné staré statické sitemap/robots konfigurace nahraďte.
- Open Graph metadata jsou přímo v serverovém HTML podle [Open Graph](https://ogp.me/),
  strukturovaná data používají [NewsArticle](https://schema.org/NewsArticle).
  Po nasazení ověřte veřejnou URL a dostupnost obrázku ve Facebook Sharing
  Debuggeru. Facebook může starší náhled cachovat; po změně použijte nové načtení.
- Další kategorie lze přidávat SQL migrací; redakční správu kategorií MVP neobsahuje.

### AI asistent článku

AI asistent je volitelný a funguje pouze v autentizované administraci. Do pole
**Podklady pro AI** vložte původní text, tiskovou zprávu nebo vlastní ověřené
poznámky. Samotná URL nestačí; aplikace webové stránky automaticky nestahuje.
Volitelně přidejte redakční pokyn. AI připraví titulek, perex, obsah v
Markdownu, kategorii, slug a SEO pole. Výsledek pouze vyplní formulář, nic
neukládá a u nového článku vždy ponechá stav konceptu.

Do serverového `.env` přidejte:

```dotenv
OPENAI_API_KEY='sk-...'
OPENAI_NEWS_MODEL=gpt-6.1-sol
```

API klíč musí zůstat pouze na serveru a nesmí se posílat do prohlížeče ani
commitovat. `OPENAI_NEWS_MODEL` je volitelný; výchozí je `gpt-6.1-sol`.
Po změně konfigurace aplikaci restartujte. Bez API klíče je tlačítko AI
vypnuté a ruční editor dál funguje.

Integrace volá serverově OpenAI Responses API se Structured Outputs a
`store: false`. Podklady a redakční pokyn se při každém generování odešlou do
OpenAI API. Neposílejte proto neveřejné osobní nebo důvěrné údaje. Volání má
časový limit, omezení frekvence a může být účtované podle nastavení OpenAI
projektu. Model dostává instrukci nevymýšlet fakta a vrací poznámky k ověření,
ale odpovědnost za kontrolu jmen, dat, čísel, zdroje a práv k obsahu zůstává na
redaktorovi.

### Testy zpráv

`go test ./...` ověřuje validaci, autentizaci, CSRF, náhled bez zápisu,
bezpečný Markdown, strukturovanou AI odpověď, metadata, escapování a HTTP
routy. Testy AI používají lokální falešný HTTP server a nespotřebovávají API
kredit. Pro databázové integrační testy nastavte
`NEWS_TEST_DSN` na **testovací MariaDB** (účet musí smět vytvářet a mazat databáze):

```bash
NEWS_TEST_DSN='user:password@tcp(127.0.0.1:3306)/?parseTime=true' go test ./internal/store -run TestNewsStoreIntegration -v
```

Test vytvoří vlastní databázi `tenelife_news_test_*`, ověří migraci,
publikaci/stažení, filtr, stránkování, konflikty a rollback a databázi odstraní.
Bez této proměnné se databázový test přeskočí.
