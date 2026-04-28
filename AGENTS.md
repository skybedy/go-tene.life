# AGENTS.md

## Role Codexu v projektu

Codex je technicky opatrny spolupracovnik pro udrzbu a rozvoj projektu `go-tene.life`. Ma umet navazat pouze z aktualniho stavu repozitare a souboru `AGENTS.md`, `PROJECT_CONTEXT.md`, `TODO.md` a `DECISIONS.md`.

Codex nesmi predpokladat kontext ze starsich dlouhych chatu. Pred praci si vzdy precte aktualni soubory, overi stav repozitare a az potom navrhuje nebo dela zmeny.

## Zjisteny stack projektu

- Backend: Go 1.24, Echo v4.
- Databaze: MariaDB/MySQL pres `github.com/go-sql-driver/mysql`.
- Konfigurace: `.env` nacitana pres `github.com/joho/godotenv`.
- Frontend: server-side HTML sablony v `views/`, Vanilla JavaScript v `public/js`.
- Stylovani: Tailwind CSS 4 pres `@tailwindcss/cli`, vstup `resources/css/app.css`, vystup `public/css/app.css`.
- Migrace: SQL migrace v `migrations/`, spoustene pres `golang-migrate` v `Makefile`.
- Staticka data/cache: `data/*.json`, v gitu ignorovano jako dynamicka data.

## Obecne preference majitele projektu

- Hlavni jazyk preferuj Go, pokud projekt neurcuje jinak.
- Frontend preferuj Vanilla JavaScript.
- UI navrhuj jednoduse, ciste a prakticky.
- Pokud je potreba stylovani, preferuj Tailwind.
- Nepridavej zbytecne slozite frameworky.
- Vyvojove prostredi je Linux Mint.
- Server byva Ubuntu VPS.

## Pravidla prace

- Nejdřív vždy čti aktuální stav projektu.
- Vzdy zkontroluj `git status`.
- Nepredpokladej kontext ze starych chatu.
- Pred vetsi zmenou strucne popis plan.
- Po zmene spust dostupne testy nebo build.
- Dulezita rozhodnuti zapisuj do `DECISIONS.md`.
- Aktualni stav zapisuj do `PROJECT_CONTEXT.md`.
- Dalsi kroky zapisuj do `TODO.md`.
- Nepridavej do commitu `.env` ani jine citlive soubory.
- Nesahej na nesouvisejici zmeny v pracovnim stromu, pokud nejsou nutne pro aktualni ukol.
- Pri upravach CSS preferuj zdrojovy soubor `resources/css/app.css` a nasledne regeneruj `public/css/app.css` pres npm build.
- Pri databazovych zmenach pridej migraci do `migrations/` a aktualizuj `db/schema.sql`, pokud je to relevantni.

## Obvykle prikazy

- Stav repozitare: `git status --short`.
- Spusteni aplikace: `go run .` po nastaveni `.env` a dostupne MariaDB.
- Go testy: `go test ./...`.
- Go build: `go build ./...`.
- CSS build: `npm run build`.
- CSS watch: `npm run dev`.
- Migrace nahoru: `make migrate-up`.
- Migrace dolu o jednu: `make migrate-down`.
- Stav migraci: `make migrate-status`.
- Export schema snapshotu: `make dump-schema`.
