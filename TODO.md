# TODO.md

## Aktuální úkoly

- Udržovat AI kontextové soubory (`AGENTS.md`, `PROJECT_CONTEXT.md`, `TODO.md`, `DECISIONS.md`) aktuální při každé významné změně.
- Dokončeno: Zpracování všech lekcí i velkých audio souborů (1–250, 251–500) do časovaných segmentů.
- Po ověření přepnout klienty ingestu na posílání `measured_at` (datum+čas) pro ruční teplotu moře.
- Dokončeno: Statistiky s agregacemi moře (`weekly/monthly`) používají hodnoty dopočítané z `water_temperatures` přes weather processor a agregační tabulky.

## K doplnění

- Doplnit stručný popis produkčního deploy postupu (pokud se používá).
- Upřesnit provozní režim schedulerů kolektorů v produkci (jen in-app vs. kombinace s cron).
- Doplnit mapu API endpointů do samostatné dokumentace (pokud chybí).

## K ověření

- Ověřit, zda je `README.md` plně aktuální vůči současným routám a stránkám.
- Ověřit aktuální CI/CD postup (zatím nezjištěno).
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

## Možné budoucí úpravy

- Přidat stručný troubleshooting section pro lokální spuštění (DB, `.env`, migrace).
- Přidat jednotné release/checklist instrukce před nasazením.
- Ověřit nasazení migrace `0011_drop_weather_daily_sea_temperature` na všech prostředích.

- Copernicus viewer: doplnit archiv měsíců (výběr období) a produkční data exportovaná z Meteodata.
- Copernicus viewer: nahradit dočasné demo snímky reálnými SST framy.
