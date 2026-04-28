# TODO.md

## Teď

- [ ] Projit prvni verzi stranky `/spanelsko-ceska-slovicka` v prohlizeci a doladit nazvy lekci / poradi.
- [ ] Zkontrolovat, jestli README obsahuje kompletni prikazy pro spusteni aplikace od nuly.
- [ ] Rozhodnout, jestli maji byt runtime ukazkova data v `data/` dokumentovana jako potrebna seed/cache data.

## Další kroky

- [ ] Doplnit README.md o rychly local setup vcetne `.env.example`, migraci a spusteni serveru.
- [ ] Pridat zakladni kontrolu dostupnosti hlavni stranky nebo API endpointu.
- [ ] Zkontrolovat, zda je potreba systemd unit nebo deploy dokumentace pro Ubuntu VPS.
- [ ] Sjednotit, jestli ma byt commitovany CSS vystup minifikovany nebo citelny neminifikovany.
- [ ] Udrzovat `PROJECT_CONTEXT.md` po vetsich zmenach aktualni.
- [ ] Zvážit manifest pro `public/sounds/`, pokud bude potreba rucni nazev, popis, kategorie nebo vlastni razeni zvukovych lekci.

## Později

- [ ] Zvážit Dockerfile nebo docker-compose pro lokalni vyvoj, pokud to zacne setrit cas.
- [ ] Rozsirit testy pro web handlery a collectory.
- [ ] Zvážit dokumentaci monitoringu collectorů a alertingu.
- [ ] Zvážit pokrocile funkce pro slovicka: oznaceni poslechnuto, opakovani lekce, rychlost prehravani nebo textovy prepis.
