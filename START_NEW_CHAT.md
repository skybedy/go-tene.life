# Start nového chatu s AI asistentem

Pracuj v aktuálním repozitáři projektu.

Nepředpokládej žádný kontext z předchozího chatu. Navazuj pouze na aktuální stav souborů, dokumentace a pracovního stromu.

Nejdřív si přečti dostupné kontextové soubory, pokud existují:

- `AGENTS.md`
- `PROJECT_CONTEXT.md`
- `TODO.md`
- `DECISIONS.md`
- `START_NEW_CHAT.md`
- `README.md`

Pokud některý soubor neexistuje, nevadí — pokračuj bez něj.

Potom zkontroluj stav repozitáře:

```bash
pwd
git branch --show-current
git status --short