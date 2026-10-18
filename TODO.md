# TODO.md

## Aktuální úkoly

- Udržovat AI kontextové soubory (`AGENTS.md`, `PROJECT_CONTEXT.md`, `TODO.md`, `DECISIONS.md`) aktuální při každé významné změně.
- Dokončeno: Zpracování všech lekcí i velkých audio souborů (1–250, 251–500) do časovaných segmentů.
- Po ověření přepnout klienty ingestu na posílání `measured_at` (datum+čas) pro ruční teplotu moře.
- Dokončeno: Statistiky s agregacemi moře (`weekly/monthly`) používají hodnoty dopočítané z `water_temperatures` přes weather processor a agregační tabulky.

## K doplnění

- Upřesnit provozní režim schedulerů kolektorů v produkci (jen in-app vs. kombinace s cron).
- Doplnit mapu API endpointů do samostatné dokumentace (pokud chybí).

## K ověření

- Vizuálně ověřit nový `/statistics/monthly-daily` (select měsíce, tabulka po dnech, grafy) na mobilu i desktopu, ve všech jazykových verzích.
- Ověřit, zda je `README.md` plně aktuální vůči současným routám a stránkám.
- Ověřit, zda ruční produkční deploy přes `deploy.sh` nemá být nahrazen bezpečnějším release/CI postupem.
- Ověřit kontrolní SQL dotazy z `docs/sql/water_temperatures_migration_checks.sql` po backfillu.
- Ověřit vizuálně weather box na mobilu i desktopu po swapu barev (kontrast, čitelnost).
- Vizuálně ověřit na `/statistics/daily` chování teploty moře po restartu aplikace: `Dnes HH:MM` pro dnešní měření a fallback na poslední dostupné datum+čas, pokud dnešní záznam chybí.
- Vizuálně ověřit na `/statistics/monthly`, `/statistics/weekly` a `/statistics/annual`, že spodní sea graf zobrazuje agregační hodnoty `sea_temperature` z API `/api/weather/{type}` (ne syrovou historii měření).
- Vizuálně ověřit na `/statistics/monthly`, že horní karty „aktuální měsíc průběžně k poslednímu datu“ odpovídají dynamickému výpočtu z raw `weather` a `water_temperatures` oříznutému podle posledního dne ve `weather_daily`.
- Vizuálně ověřit na `/statistics/recent` čitelnost nově zobrazeného času měření moře (formát `H:MM` bez počáteční nuly) na mobilu i desktopu.
- Vizuálně ověřit na `/spanelsko-ceska-slovicka`, že volba `Přehrát všechno za sebou` používá nový track `1-500` a synchronizuje panel slovíček bez posunu.
- Vizuálně ověřit nové routy a podmenu vizualizací teplot Atlantiku:
  - `/atlantic-sea-temperature-tenerife`
  - `/atlantic-sea-temperature-canary-islands`
  včetně locale prefix variant (`/en/...`, `/es/...`, ...).
- Vizuálně ověřit desktop/mobil navigaci po přidání nové položky a podmenu (zalomení kolem `lg` breakpointu).
- Vizuálně ověřit lokalizované texty na obou stránkách teploty moře a ve viewer chybových hláškách.
- Vizuálně ověřit nový tmavě modrý styl homepage weather boxu a tlačítka „Zvětšit fotku“ na mobilu i desktopu (mezery, rámečky, kontrast, zaoblení).
- Vizuálně ověřit čitelnost textu a viditelnost fotografie pod průhlednějším weather info boxem při šířkách do 640 px.
- Po nasazení ověřit na mobilu, že URL `/css/app.css?v=…` dostala nový otisk a nový vzhled se načte bez ručního vymazání cache.
- Vizuálně ověřit vertikální centrování přílivu a odlivu pro dny s jednou i dvěma událostmi, včetně času a výšky na úzkém mobilu.
- Po nasazení vizuálně ověřit redesign homepage z původní experimentální větve na mobilu i desktopu.

## Možné budoucí úpravy

- Přidat stručný troubleshooting section pro lokální spuštění (DB, `.env`, migrace).
- Přidat jednotné release/checklist instrukce před nasazením.
- Ověřit nasazení migrace `0011_drop_weather_daily_sea_temperature` na všech prostředích.

- Copernicus viewer: doplnit archiv měsíců (výběr období) a produkční data exportovaná z Meteodata.
- Copernicus viewer: nahradit dočasné demo snímky reálnými SST framy.

## Zprávy z Tenerife (2026-10-09)

- Dokončeno: databázová vrstva a migrace 0012, veřejné stránky, kategorie, stránkování, administrace, náhled, publikace/stažení, navigace, homepage, OG/SEO/NewsArticle, sitemap a testy.
- Dokončeno: Markdown editor s formátovací lištou a bezpečným serverovým náhledem; volitelný AI asistent generující všechna redakční pole jako neuložený koncept přes strukturovanou odpověď OpenAI.
- Před použitím AI nastavit serverový `OPENAI_API_KEY`, zkontrolovat rozpočet/limity OpenAI projektu a provést redakční zkoušku na reálných veřejných podkladech.
- Před nasazením ověřit DB verzi, aplikovat 0012, nastavit SITE_URL a redakční přihlašovací údaje, sestavit CSS/binárku a restartovat. Podrobný postup v README.
- Po nasazení ověřit Nginx HTTPS/cache a průchod sitemap/robots do aplikace, poté skutečný Facebook náhled první publikované URL.
- Po nasazení vizuálně ověřit lištu Markdown editoru a AI generování v desktopovém i mobilním formuláři.
- Samostatně opravit existující migraci 0006 pro instalace od nuly (kolize FK po 0005 a očekávání legacy station_id).
- Samostatně opravit instalaci od nuly také u migrace 0007 (duplicitní sloupce avg_sea_temperature z 0001).
- Vizuální kontrolu mobilu/desktopu a formuláře dokončit v dostupném prohlížeči; v této relaci nástroj nehlásil žádný dostupný prohlížeč. HTTP smoke test stávajících i nových stránek prošel.
- Dokončeno: test publikačního cyklu přes autentizované HTTP proti izolované MariaDB, kontrola sitemap/homepage po stažení, `go test -race ./...`, `go vet`, Go/CSS build.
