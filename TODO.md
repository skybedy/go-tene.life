# TODO.md

## Aktuální úkoly

- Udržovat AI kontextové soubory (`AGENTS.md`, `PROJECT_CONTEXT.md`, `TODO.md`, `DECISIONS.md`) aktuální při každé významné změně.
- Otestovat migrace `0008` a `0009` na lokální kopii produkčních dat.
- Po ověření přepnout klienty ingestu na posílání `measured_at` (datum+čas) pro ruční teplotu moře.

## K doplnění

- Doplnit stručný popis produkčního deploy postupu (pokud se používá).
- Upřesnit provozní režim schedulerů kolektorů v produkci (jen in-app vs. kombinace s cron).
- Doplnit mapu API endpointů do samostatné dokumentace (pokud chybí).

## K ověření

- Ověřit, zda je `README.md` plně aktuální vůči současným routám a stránkám.
- Ověřit aktuální CI/CD postup (zatím nezjištěno).
- Ověřit kontrolní SQL dotazy z `docs/sql/water_temperatures_migration_checks.sql` po backfillu.

## Možné budoucí úpravy

- Přidat stručný troubleshooting section pro lokální spuštění (DB, `.env`, migrace).
- Přidat jednotné release/checklist instrukce před nasazením.
- Po potvrzení kompletní migrace připravit ostrou cleanup migraci pro drop `weather_daily.sea_temperature` (zatím jen šablona v `migrations/cleanup_later_drop_weather_daily_sea_temperature.sql`).
