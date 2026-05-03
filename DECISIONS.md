# DECISIONS.md

## 2026-04-28

### Dulezita technicka rozhodnuti

- Projekt zustava primarne Go aplikace s Echo v4.
- Frontend zustava jednoduchy: HTML sablony, Vanilla JavaScript a Tailwind CSS.
- Databazove zmeny se maji delat pres verzovane SQL migrace v `migrations/`.
- `.env` a jine citlive soubory se necommituji.
- Kontext pro dalsi Codex chaty se udrzuje primo v repozitari v souborech `AGENTS.md`, `PROJECT_CONTEXT.md`, `TODO.md` a `DECISIONS.md`.

### Pouzite technologie

- Go 1.24 / Echo v4.
- MariaDB/MySQL.
- Tailwind CSS 4.
- Vanilla JavaScript.
- SQL migrace pres `golang-migrate`.

### Duvody dulezitych voleb

- Go odpovida aktualnimu projektu a preferenci majitele.
- Vanilla JavaScript a server-side sablony udrzuji frontend jednoduchy bez zbytecnych frameworku.
- Tailwind je uz v projektu zavedeny a je vhodny pro rychle prakticke UI upravy.
- Migrace a schema snapshot zlepsuji opakovatelnost databazovych zmen.
- Repo-local kontext brani zavislosti projektu na historii jednoho dlouheho Codex chatu.

### 2026-04-28: Stranka se zvukovymi slovicky

- Pro spanelsko-ceske zvukove lekce vznikla samostatna stranka `/spanelsko-ceska-slovicka`.
- MP3 soubory zustavaji v `public/sounds/`, aby byly soucasti verejnych statickych souboru projektu.
- Verejne URL pro samotne audio soubory pouziva prefix `/spanelsko-ceska-slovicka/files/`, aby URL odpovidala obsahu stranky.
- Stara routa `/sounds` zustava jako presmerovani na novou adresu.
- Playlist se zatim generuje z dostupnych `.mp3` souboru podle nazvu souboru bez databaze; to je jednoduche a vhodne pro prvni iteraci.
- Souhrnne soubory `spanelsko_ceska_slovicka_1_250.mp3` a `spanelsko_ceska_slovicka_251_500.mp3` zustavaji jako horni prehledove volby; nezobrazuji se dole mezi jednotlivymi lekcemi.
- Horni volby nemaji mit modre primarni zvyrazneni; pouziva se neutralni vzhled a modra zustava jen jako hover/focus signal.

### Otevrene otazky

- Neni zatim rozhodnuto, jestli projekt potrebuje Docker konfiguraci.
- Neni zatim rozhodnuto, jak presne dokumentovat produkcni deploy na Ubuntu VPS.
- Neni zatim rozhodnuto, jestli `data/` ma zustat ciste runtime cache, nebo zda nektera ukazkova data patri do repozitare.
- Neni zatim rozhodnuto, zda ma byt `public/css/app.css` v repozitari udrzovany minifikovany pres `npm run build`, nebo citelny neminifikovany.
- Neni zatim rozhodnuto, jestli budou zvukove lekce pozdeji potrebovat metadata v databazi nebo rucne udrzovany manifest.
