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

10. Denní statistiky používají pro teplotu moře denní point-in-time hodnotu z `water_temperatures` místo `weather_daily.sea_temperature`.
- Důvod: zdroj `water_temperatures` obsahuje přesné timestampy měření a umožňuje ve statistice zobrazit i datum+čas poslední aktualizace.

11. Legacy sloupec `weather_daily.sea_temperature` se odstraňuje migrací `0011_drop_weather_daily_sea_temperature`.
- Důvod: po přepnutí čtení i zápisu na `water_temperatures` už sloupec způsoboval jen duplicitní stav a riziko nekonzistence.

12. Na `/statistics/daily` se při chybějícím dnešním měření teploty moře použije fallback na poslední dostupné měření, ale štítek `Dnes HH:MM` se zobrazí pouze pro skutečně dnešní záznam.
- Důvod: UI má korektně rozlišit aktuální dnešní měření od historické fallback hodnoty a předejít dojmu, že starší data jsou dnešní. Čas fallback záznamu se pro tuto stránku vrací tak, jak je uložený ve `water_temperatures`.

13. Pro volbu `Přehrát všechno za sebou` na `/spanelsko-ceska-slovicka` se používá dedikovaný spojený track `spanelsko_ceska_slovicka_1_500.mp3` s vlastním seznamem časovaných segmentů.
- Důvod: playlistové řazení více samostatných souborů negarantuje jednu souvislou časovou osu segmentů; spojený soubor umožní přesné synchronní zobrazování CZ/ES slovíček i v režimu „všechno za sebou“.

14. Na stránkách `/statistics/weekly`, `/statistics/monthly` a `/statistics/annual` se sea graf plní agregačním datasetem `sea_temperature` z endpointu `/api/weather/{type}`, nikoli syrovou historií měření.
- Důvod: grafulárně má zobrazovat stejnou agregační úroveň jako tabulka dané stránky (týden/měsíc), ne point-in-time křivku z `/api/water-temperatures/history`, která vizuálně neodpovídá měsíčním/týdenním průměrům.

15. Na `/statistics/recent` se u hodnoty moře zobrazuje i čas měření z `SeaMeasuredAt` ve formátu `H:MM` bez úvodní nuly.
- Důvod: uživatel potřebuje vidět, kdy byla denní sea hodnota skutečně naměřena; kratší formát je čitelnější v úzké tabulce.

16. Na `/statistics/recent` se měrné jednotky zobrazují v hlavičkách sloupců, ne v každé datové buňce.
- Důvod: odstranění vizuálního šumu v tabulce; hodnoty jsou lépe porovnatelné po řádcích.

17. V grafových osách statistik se hodnoty formátují na 1 desetinné místo před zobrazením.
- Důvod: eliminace floating-point artefaktů v UI (např. dlouhé desetinné zbytky u teplot).

18. V homepage hodinové tabulce se měrné jednotky zobrazují v hlavičkách sloupců a buňky obsahují pouze hodnoty.
- Důvod: sjednocení s tabulkami statistik a lepší čitelnost opakovaných hodinových hodnot.

19. Souhrnné karty na `/statistics/recent` používají stejnou vizuální strukturu jako karty na `/statistics/daily`.
- Důvod: sjednocení statistik a odstranění odlišného měřítka typografie mezi podobnými stránkami.

20. Průměrná teplota moře na `/statistics/recent` se počítá klientsky z dostupných denních hodnot `datasets.sea_temperature`.
- Důvod: endpoint už vrací jednu denní hodnotu moře z `water_temperatures`; souhrnná karta má zobrazovat průměr za stejný rozsah jako ostatní recent metriky.

21. Souhrnné karty na `/statistics/monthly` zobrazují aktuální měsíc průběžně k poslednímu datu z dynamického endpointu `/api/weather/monthly-current`, ne z denních ani měsíčních agregačních tabulek.
- Důvod: horní souhrn má odpovídat otevřenému měsíci, ale nemá ukazovat neuzavřený dnešek. Poslední datum se proto bere z `weather_daily`; pokud `water_temperatures` nemá měření ve stejný den, používá se stále stejné datum z `weather_daily` a moře se průměruje jen z dostupných měření v tomto období.

22. Tabulky na homepage a statistikách používají jednotný styl podle `/statistics/monthly` se střídavým podbarvením řádků.
- Důvod: sjednocená typografie, odsazení, centrování a `odd/even` řádky zvyšují čitelnost napříč podobnými datovými tabulkami.

- 2026-05-09: Pro první fázi Copernicus integrace v tene.life zvolena statická architektura (`public/data/.../manifest.json` + snímky), bez volání externích API z webu.
- 2026-05-09: Přidána route `/teplota-more` a samostatný vanilla JS přehrávač s Play/Pause + sliderem, aby šlo později jednoduše přidat archiv měsíců.
- 2026-05-10: Pro veřejné odkazy na viewer zvoleny anglické a jednoznačné cesty:
  - `/atlantic-sea-temperature-tenerife`
  - `/atlantic-sea-temperature-canary-islands`
  Staré české routy (`/teplota-more`, `/teplota-more-kanary`) zůstávají jako redirecty kvůli kompatibilitě.
- 2026-05-10: Položka „Vizualizace teplot Atlantiku“ je v hlavní navigaci jako samostatné podmenu (`Tenerife`, `Kanárské ostrovy`), nikoli pod „Rozšířené statistiky počasí`.
- 2026-05-10: Status text v vieweru se běžně nezobrazuje; je vyhrazen jen pro chybové stavy (např. chybějící manifest / chyba načtení snímku), aby UI nebylo zbytečně zahlcené.
- 2026-05-10: Nové texty pro stránky teploty moře a viewer JS hlášky jsou řešené přes i18n klíče pro všechny podporované jazyky (`cs`, `en`, `es`, `pl`, `de`, `fr`, `it`, `hu`).
- 2026-05-10: Homepage weather info box a tlačítko „Zvětšit fotku“ používají sjednocený tmavě modrý vizuální styl s modrým rámem, jemnou světlou linkou a stejným vnějším paddingem.
- Důvod: uživatelský požadavek sjednotit vzhled podle dodané reference a zachovat stávající layout, velikosti a obsah.
