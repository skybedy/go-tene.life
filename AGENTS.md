# AGENTS.md

Pravidla pro AI coding agenty v tomto repozitáři (`go-tene.life`).

## Základní principy

- Komunikuj primárně v češtině, ale angličtina je akceptovatelná v technických kontextech.
- Pracuj pouze podle aktuálního stavu souborů a repozitáře.
- Nepředpokládej kontext z předchozích chatů.
- Před změnami vždy ověř aktuální adresář (`pwd`) a Git stav, pokud je Git dostupný.
- Před úpravami si načti relevantní soubory, kterých se změna týká.

## Práce s Gitem a změnami

- Nevracej (`revert`) nesouvisející necommitované změny.
- Nepřepisuj ruční práci uživatele.
- Nedělej destruktivní Git operace bez výslovného pokynu.
- Dělej malé, cílené změny s jasným účelem.

## Technický styl projektu

- U existujícího projektu nedělej velké refaktory bez výslovného důvodu.
- Zachovej stávající technologie, architekturu a styl projektu.
- Nepřidávej nové závislosti bez jasného důvodu.
- Neměň aplikační kód mimo rozsah aktuálního úkolu.

## Ověření po změnách

- Pokud je to možné a bezpečné, spusť relevantní testy nebo build.
- Když ověření není možné (např. chybí služby/credentials), napiš to explicitně.

## Údržba AI kontextu

- Pokud se změnil stav projektu, aktualizuj:
  - `PROJECT_CONTEXT.md`
  - `TODO.md`
  - `DECISIONS.md`
- Udržuj zápisy stručné, faktické a ověřené.
