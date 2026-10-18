# PROJECT_CONTEXT.md

## Projekt

- Název: `go-tene.life`
- Účel: webová aplikace pro informace související s Tenerife (počasí, vlny, kvalita vody, příliv/odliv, PWS teploty).

## Aktuální stav

- Typ: rozpracovaný existující projekt.
- Git: redesign homepage, stránka denních statistik po měsících a MVP zpráv jsou sloučené do `main`.
- Všechna slovíčka (lekce 1–26 i velké soubory 1–250, 251–500) mají vygenerované časované segmenty v `public/sounds/index.json`.

## Technologie

- Backend: Go `1.24`, Echo v4.
- Databáze: MariaDB/MySQL přes `github.com/go-sql-driver/mysql`.
- Konfigurace: `.env` přes `github.com/joho/godotenv`.
- Frontend: server-side HTML šablony (`views/`) + Vanilla JS (`public/js`).
- Styling: Tailwind CSS 4 (`@tailwindcss/cli`), zdroj `resources/css/app.css`, výstup `public/css/app.css`.
- URL hlavního CSS obsahuje automatický krátký SHA-256 otisk vestavěného `app.css`; po změně CSS se tedy změní URL a prohlížeče nepoužijí zastaralou cache.
- Migrace: SQL migrace v `migrations/` přes `golang-migrate` (Makefile).
- Sea temperature redesign (active): nové měření jde do `water_temperatures` (timestamp + hodnota), legacy sloupec `weather_daily.sea_temperature` je určen k odstranění migrací `0011`.
- Backfill pravidlo: historické hodnoty z `weather_daily.sea_temperature` se mapují na čas `10:00:00` (UTC) jako default pro legacy denní záznamy.
- Denní statistiky (`/statistics/daily`) už berou teplotu moře z `water_temperatures` (nejnovější měření v rámci dne) a v UI zobrazují i timestamp poslední aktualizace.
- Na `/statistics/daily` je doplněný fallback: pokud pro dnešní datum chybí teplota moře, API vrátí poslední dostupné měření; UI pak zobrazuje `Dnes HH:MM` jen pokud měření opravdu pochází z dnešního dne, jinak zobrazí celé datum+čas uložené v `water_temperatures`.
- Na `/statistics/monthly` a dalších agregačních stránkách byl opraven sea graf: už nebere syrovou historii z `/api/water-temperatures/history`, ale agregační dataset `sea_temperature` z `/api/weather/{weekly|monthly|annual}`.
- Na `/statistics/recent` se ve sloupci moře zobrazuje i čas měření (`SeaMeasuredAt`) ve formátu `H:MM` bez počáteční nuly.
- Na `/statistics/recent` byly jednotky přesunuty z hodnot v buňkách do hlaviček sloupců (°C, hPa, %) pro přehlednější tabulku.
- V grafech statistik je osa Y nově zaokrouhlená na 1 desetinné místo, aby se neukazovaly floating-point artefakty (např. `21.400000000000006`).
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
- Produkční deployment běží ručně na VPS skriptem `deploy.sh`: aktualizuje aktuální větev z GitHubu, sestaví Tailwind a Go binárku `tenelife-app`, odstraní zdrojové soubory a následně vyžaduje `sudo systemctl restart tenelife`. Databázové migrace skript nespouští.
- V homepage weather boxu proběhl swap barev (to, co bylo `text-white/90`, je nyní `text-orange-300`; původní `text-orange-300` je nyní `text-white/90`).
- Homepage hodinová tabulka má jednotky v hlavičkách sloupců (`°C`, `hPa`, `%`) a hodnoty v buňkách bez jednotek, se zarovnáním na střed.
- Homepage weather info box a tlačítko „Zvětšit fotku“ mají do šířky 640 px zcela průhledný zadní layoutový obal bez rámečku, stínu a rozmazání; viditelné zůstávají jen jemně tónované vnitřní dlaždice/tlačítko, aby fotografie webkamery nebyla zakrytá. Větší breakpointy zachovávají tmavý vzhled. Další vizuální změny informačního boxu se mají vždy současně promítnout i do tlačítka.
- Sloupce přílivu a odlivu drží nadpis nahoře a ve zbývající výšce vertikálně centrují jednu nebo dvě události; každá událost vždy zobrazuje čas i výšku moře. Do 420 px dostávají více prostoru na úkor sloupce moře a hodnoty se nezalamují.
- Tabulky na homepage, `/statistics/daily`, `/statistics/weekly`, `/statistics/monthly` a `/statistics/recent` používají sjednocený monthly styl: centrované buňky, jednotky v hlavičkách, `text-xs sm:text-sm`, `px-3 py-3` a střídavé podbarvení řádků.
- Horní tři souhrnné karty na `/statistics/recent` používají stejný formát a typografii jako souhrnné karty na `/statistics/daily`.
- `/statistics/recent` má doplněnou čtvrtou souhrnnou kartu pro průměrnou teplotu moře z dostupných hodnot `datasets.sea_temperature`.
- Horní souhrn na `/statistics/monthly` zobrazuje „aktuální měsíc průběžně k poslednímu datu“ přes `/api/weather/monthly-current`; poslední datum se bere z posledního záznamu `weather_daily` v aktuálním měsíci a výpočet běží z raw `weather` a `water_temperatures` jen do konce tohoto dne.

- 2026-05-09: Přidána první statická stránka `/teplota-more` pro Copernicus SST vizualizace. Načítá manifest z `public/data/copernicus/.../manifest.json` a používá nový JS viewer (`public/js/copernicus-viewer.js`).
- 2026-05-10: Viewer pro teplotu moře podporuje výběr období (měsíce) ze statického seznamu a načítá měsíční manifesty z `public/data/copernicus/sea-temp/{tenerife|canary}/YYYY/MM/manifest.json`.
- 2026-05-10: Přidána druhá stránka pro oblast Kanárských ostrovů a nové anglické routy:
  - `/atlantic-sea-temperature-tenerife`
  - `/atlantic-sea-temperature-canary-islands`
  Staré `/teplota-more` a `/teplota-more-kanary` zůstaly jako redirecty kvůli zpětné kompatibilitě.
- 2026-05-10: V navigaci je nová položka s podmenu „Vizualizace teplot Atlantiku“ (`Tenerife`, `Kanárské ostrovy`) v desktop i mobilní variantě.
- 2026-05-10: Texty a UI pro obě stránky teploty moře (včetně hlášek vieweru v JS) jsou lokalizované pro všechny podporované jazyky (`cs`, `en`, `es`, `pl`, `de`, `fr`, `it`, `hu`).
- 2026-05-10: Homepage weather info box má upravený tmavě modrý vizuální styl podle dodané reference. Stejný rám, pozadí a vnější padding používá i tlačítko „Zvětšit fotku“; tlačítko má bílé písmo/ikonu.
- 2026-05-12: UX/UI redesign homepage sjednotil panelové kontejnery (`surface-panel`) a jemně navýšil čitelnost weather overlaye (větší typografie a vyšší první řádek gridu).
- 2026-05-12: Navigace byla přepojena na `public/images/tenelife-logo-transparent-v2.png`, protože původní `tenelife-logo.png` obsahuje neprůhledné pozadí `rgb(243,244,246)`.
- 2026-06-28: Přidána nová stránka `/statistics/monthly-daily` („Dny v měsíci“) v navigaci statistik vedle Denní/Týdenní/Měsíční/Posledních 10 dní. Zobrazuje denní rozpis vybraného měsíce ve stylu tabulky `/statistics/recent`. Výběr měsíce je `<select>` naplněný dynamicky podle reálně existujících dat (`SELECT DISTINCT YEAR(date), MONTH(date) FROM weather_daily`), takže obsahuje jen měsíce od listopadu 2025 dál. Pokud aktuální měsíc ještě nemá data, automaticky se vybere nejnovější dostupný měsíc. Nový handler `MonthlyDailyStatisticsHandler` a store metoda `GetAvailableMonths()`; grafy znovu používají existující API `/api/weather/monthly-daily?year=&month=`.

- 2026-10-09: Implementováno MVP českých zpráv: `/zpravy`, `/zpravy/:slug`, `/admin/zpravy`, `/sitemap.xml`, `/robots.txt`. Sdílí Echo, MySQL pool, Go šablony a Tailwind; homepage zobrazuje poslední 3 publikované zprávy. Kategorie a články přidává nedestruktivní up migrace `0012_news`.
- Redakce používá jeden účet (`NEWS_ADMIN_USERNAME`, bcrypt `NEWS_ADMIN_PASSWORD_HASH`), helper `go run ./cmd/news-password`, Basic Auth za HTTPS, omezení přihlášení a CSRF. Obě prázdné hodnoty administraci vypnou. `SITE_URL` řídí canonical a absolutní OG URL.
- Helper pro redakční heslo přijímá 6–72 bajtů a vytváří bcrypt hash s costem 12; administrace dále omezuje frekvenci pokusů o přihlášení.
- Editor používá bezpečně renderovaný Markdown s lištou pro běžné formátování, náhled bez zápisu, koncept/publikaci/stažení a potvrzené mazání. Raw HTML a nebezpečné odkazy Goldmark nevykreslí. První publikace uzamkne slug a datum; revize chrání před přepsáním souběžné změny. Obrázky se zadávají URL (HTTPS nebo `/images/…`), soubory se nahrávají mimo editor.
- Volitelný AI asistent v administraci používá serverově OpenAI Responses API se Structured Outputs. Z vložených podkladů připraví titulek, perex, Markdown obsah, kategorii, slug a SEO pole; nic automaticky neukládá ani nepublikuje. Konfigurace: `OPENAI_API_KEY`, volitelně `OPENAI_NEWS_MODEL` (výchozí `gpt-6.1-sol`). Podklady se neposílají nikam, pokud redaktor nestiskne tlačítko generování.
- Integrační testy na izolované MariaDB 10.11 ověřily migraci 0012, publikaci/stažení, trvalý slug, konflikty, filtrování, stránkování a rollback. `NEWS_TEST_DSN` umožňuje opakování; bez proměnné se DB test přeskočí. Žádné migrace nebyly aplikovány na existující databázi.
- Při sestavení čerstvé testovací DB zjištěn existující problém migrace 0006: po 0005 koliduje název FK `pws_latest_station_ref_id_fk`; migrace navíc očekává starší sloupec `station_id`. Netýká se přechodu 0011 → 0012; řešení starší migrace je samostatný úkol.
- Na čerstvém schématu má také migrace 0007 duplicitní `avg_sea_temperature` (již existuje v 0001). Lokální smoke test aplikace proto použil ekvivalentní schéma bez opakování 0006/0007; nová 0012 byla provedena a otestována samostatně.
- Závěrečné ověření: `go test -race ./...` včetně DB testu, `go vet ./...`, Go build a Tailwind build prošly. Přes HTTP proti izolované DB prošel celý cyklus náhled → koncept → publikace → stažení → potvrzené smazání i načtení stávajících stránek. Vizuální kontrola zůstává neověřená (žádný dostupný prohlížeč v nástroji).
- 2026-10-18: Větev `experiment/jack-dorsey-redesign` byla včetně MVP zpráv sloučena do `main`.
