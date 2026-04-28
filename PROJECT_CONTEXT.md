# PROJECT_CONTEXT.md

## Strucny popis projektu

`go-tene.life` je Go webova aplikace pro web Tene.life. Zobrazuje aktualni informace souvisejici s Tenerife: pocasi, webkameru, statistiky, teplotni mapu PWS stanic, vlny, kvalitu vody a priliv/odliv. Aplikace ma HTML sablony, API endpointy, periodicke collectory a MariaDB databazi.

## Aktualni stav

- Projekt je Go aplikace s hlavnim vstupem v `main.go`.
- Webovy server bezi na Echo v4.
- Frontend je kombinace Go HTML templates, statickych JS souboru a Tailwind CSS.
- README popisuje databazove migrace a collectory pro waves, water quality, PWS a tides.
- V repozitari nejsou nalezeny `Dockerfile` ani `docker-compose.yml`.
- Pri vytvoreni tohoto kontextu byl `git status --short` cisty.
- Dne 2026-04-28 byly overeny prikazy `go test ./...`, `GOCACHE=/tmp/go-build go build ./...` a `npm run build`.
- Dne 2026-04-28 byla v nove vetvi `feature/spanish-czech-sounds` rozpracovana samostatna stranka `/sounds` pro spanelsko-ceske zvukove lekce.
- Zvukove soubory jsou ulozene v `public/sounds/`; web je servíruje pod URL prefixem `/sounds/files/`.

## Pouzivany stack

- Go 1.24.0, toolchain `go1.24.13`.
- Echo v4.
- MariaDB/MySQL.
- `godotenv` pro `.env`.
- Tailwind CSS 4 pres npm CLI.
- Vanilla JavaScript.
- SQL migrace pres `golang-migrate`.

## Hlavni adresare a soubory

- `main.go` - vstup aplikace, konfigurace serveru, route, collectory.
- `internal/web/` - HTTP handlery a renderovani stran/API.
- `internal/store/` - prace s databazi.
- `internal/models/` - datove modely.
- `internal/pws/` - collector a klient pro Weather Company PWS API.
- `internal/tides/` - tide collector a fetch logika.
- `internal/waves/` - collector pro vlny.
- `internal/water/` - collector pro kvalitu vody.
- `internal/i18n/` - lokalizace a pomocne funkce pro locale URL.
- `internal/utils/` - validace env a pomocne funkce.
- `views/` - HTML sablony.
- `public/js/` - Vanilla JavaScript pro mapu a grafy.
- `public/sounds/` - verejne MP3 soubory pro zvukove lekce spanelsko-ceskych slovicek.
- `resources/css/app.css` - zdrojovy CSS/Tailwind vstup.
- `public/css/app.css` - vygenerovany CSS vystup.
- `migrations/` - SQL migrace.
- `db/schema.sql` - snapshot databazoveho schema.
- `data/` - runtime JSON cache; v `.gitignore` je vedena jako dynamicka data.
- `.env.example` - priklad konfigurace.
- `.env` - lokalni citliva konfigurace, nesmi do commitu.
- `Makefile` - migrace a export schema snapshotu.
- `package.json` - npm skripty pro Tailwind CSS.

## Jak projekt spustit

1. Priprav `.env` podle `.env.example`.
2. Zajisti dostupnou MariaDB/MySQL databazi.
3. Spust migrace:

```bash
make migrate-up
```

4. Spust aplikaci:

```bash
go run .
```

Port se bere z `APP_PORT`; pokud neni nastaven, aplikace pouzije `8080`.

Rucni collectory:

```bash
go run . collect:waves
go run . collect:water
go run . collect:pws
```

## Jak projekt testovat

Dostupne Go testy:

```bash
go test ./...
```

Aktualne zjistene test soubory:

- `internal/pws/wu_client_test.go`
- `internal/tides/collector_test.go`

Automatizovane frontend testy zatim nejsou zjisteny.

## Jak projekt buildit

Go build:

```bash
go build ./...
```

V sandboxovem prostredi muze byt potreba pouzit zapisovatelnou Go cache:

```bash
GOCACHE=/tmp/go-build go build ./...
```

Tailwind CSS build:

```bash
npm run build
```

Tailwind CSS watch pro vyvoj:

```bash
npm run dev
```

## Znama omezeni/problemy

- Docker konfigurace zatim neni definovana.
- `.env` obsahuje lokalni/citlive hodnoty a nesmi se commitovat.
- Spusteni aplikace a migraci vyzaduje dostupnou MariaDB/MySQL a spravne env promenne.
- `data/` je dynamicky runtime cache adresar; README ho zminuje, ale `.gitignore` ho ignoruje.
- Frontend testy zatim nejsou definovane.

## Poznamky pro dalsi navazani

- Kazdy novy Codex chat ma nejdriv precist `AGENTS.md`, tento soubor, `TODO.md` a `DECISIONS.md`.
- Pred kazdou praci over `git status --short`.
- Pri zmene UI/CSS upravuj primarne `resources/css/app.css`, sablony ve `views/` a JS v `public/js/`; pak spust `npm run build`.
- Stranka se zvukovymi slovicky je na `/sounds`; handler nacita `.mp3` soubory z `public/sounds/` pri renderovani stranky.
- Pri zmene backendu spust `go test ./...` a podle potreby `go build ./...`.
- Pri databazovych zmenach pridej migraci a aktualizuj `db/schema.sql`.
