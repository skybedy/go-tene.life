# DECISIONS.md

## Zaznamenaná technická rozhodnutí

1. Backend je postavený v Go s frameworkem Echo v4.
- Důvod: nezjištěno (vyplývá z `go.mod` a `main.go`).

2. Databázová vrstva používá MariaDB/MySQL přes `go-sql-driver/mysql`.
- Důvod: nezjištěno (vyplývá z `go.mod`, `main.go`, `Makefile`).

3. Databázové změny jsou řízené SQL migracemi přes `golang-migrate`.
- Důvod: nezjištěno (vyplývá z `migrations/` a `Makefile`).

4. Frontend je renderovaný server-side HTML šablonami.
- Důvod: nezjištěno (vyplývá z `views/` a renderingu v `main.go`).

5. Styling je řešen Tailwind CSS 4 CLI buildem z `resources/css/app.css` do `public/css/app.css`.
- Důvod: nezjištěno (vyplývá z `package.json`).

6. Externí data (waves/water/PWS) se sbírají kolektory a ukládají do cache/DB; API je čte z lokálních dat.
- Důvod: pravděpodobně stabilita výkonu a omezení závislosti na externích API během requestu; explicitní důvod zatím nezjištěn.

7. Teplota moře se migruje z denního agregátu `weather_daily.sea_temperature` do samostatné tabulky `water_temperatures` s přesným časem měření.
- Důvod: ruční měření probíhá nepravidelně a vícekrát denně; denní agregát neodpovídá realitě point-in-time měření.

8. Backfill historických teplot moře používá default čas `10:00:00` (UTC) pro záznamy převzaté z `weather_daily`.
- Důvod: v původních datech není přesný čas měření; byl zvolen jednotný a zdokumentovaný čas pro konzistentní migraci.

9. V homepage weather boxu byl použit konzistentní swap barev mezi `text-white/90` a `text-orange-300`.
- Důvod: uživatelský požadavek na sjednocení vizuální hierarchie (co bylo bílé má být oranžové a naopak).
