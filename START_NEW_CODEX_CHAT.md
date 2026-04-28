# Start nového Codex chatu

## Prompt pro Antigravity / Gemini

Pracuj v repozitari `/home/skybedy/Programming/web/go/go-tene.life`.

Nejdriv si precti:

- `AGENTS.md`
- `PROJECT_CONTEXT.md`
- `TODO.md`
- `DECISIONS.md`
- `START_NEW_CODEX_CHAT.md`

Pak spust:

```bash
git branch --show-current
git status --short
```

Nepredpokladej kontext z predchoziho chatu. Navazuj jen na aktualni stav souboru a pracovniho stromu. Nerevertuj nesouvisejici necommitovane zmeny. Pred upravami si prohledni relevantni soubory.

### Aktualni predavka

Aktualni projekt: `/home/skybedy/Programming/web/go/go-tene.life`

Aktualni vetev: `feature/spanish-czech-sounds`

Cil aktualni prace:

Doplnujeme stranku `/sounds` pro prehravani spanelsko-ceskych audio slovicek.

Dulezite zmenene/pridane soubory:

- `main.go`
- `internal/web/handlers.go`
- `internal/models/sounds.go`
- `internal/web/sounds_test.go`
- `internal/i18n/i18n.go`
- `views/nav.html`
- `views/sounds.html`
- `public/sounds/*.mp3`

Aktualni stav funkcnosti:

- Existuje route `/sounds`.
- Menu obsahuje odkaz na `Slovicka`.
- MP3 jsou v `public/sounds`.
- Audio soubory se serviruji pres `/sounds/files/...`.
- Stranka `/sounds` ma prehravac, volbu rychlosti `0.75x`, `1x`, `1.25x`, `1.5x`.
- Nahore jsou 4 zakladni volby: `Prehrat vsechno za sebou`, `Prehrat nahodne`, `Prehrat 1-250`, `Prehrat 251-500`.
- Pod tim jsou jednotlive soubory k rucnimu prehrani.
- Jednotlive soubory maji tlacitka `Prehrat` a `Prehrat ve smycce`.
- Po otevreni stranky se automaticky nevybere prvni lekce.
- Souhrnne soubory `spanelsko_ceska_slovicka_1_250.mp3` a `spanelsko_ceska_slovicka_251_500.mp3` se nemaji zobrazovat dole v seznamu jednotlivych souboru, maji zustat jen jako horni zakladni volby.
- Horni 4 volby maji byt vizualne neutralni, zadna nesmi byt modre zvyraznena jako primarni.

Pozor:

- V pracovnim stromu jsou i necommitovane kontextove soubory `AGENTS.md`, `PROJECT_CONTEXT.md`, `TODO.md`, `DECISIONS.md` a `START_NEW_CODEX_CHAT.md`. Nerevertovat.
- `public/css/app.css` je zmeneny, ale pro aktualni praci se zvuky do nej nesahej, pokud to neni nutne.
- Nemenit nazvy MP3 souboru. Zobrazovane ceske nazvy se resi v Go kodu.
- Nepredelavat frontend na framework, zustat u HTML templates, Vanilla JS a Tailwind trid.

Overeni, ktere naposledy proslo:

```bash
GOCACHE=/tmp/go-build go test ./...
GOCACHE=/tmp/go-build go build ./...
GOCACHE=/tmp/go-build go run /tmp/check_templates.go
```

Nejblizsi mozne dalsi kroky:

- Projit `/sounds` vizualne v prohlizeci.
- Doladit texty a vzhled tlacitek.
- Pripadne zvazit, jestli se pro zvuky pozdeji vyplati manifest misto generovani z nazvu souboru.

## První zpráva v novém Codex chatu

Přečti AGENTS.md, PROJECT_CONTEXT.md, TODO.md a DECISIONS.md.
Zkontroluj aktuální stav projektu přes git status.
Nepředpokládej žádný kontext ze starého chatu.

Pokračuj úkolem:
[sem napiš konkrétní úkol]

## Závěrečná zpráva na konci pracovního chatu

Aktualizuj PROJECT_CONTEXT.md, TODO.md a DECISIONS.md podle toho, co jsme právě změnili, aby šlo bezpečně navázat v novém chatu.
Potom ukaž git status a navrhni commit message.
