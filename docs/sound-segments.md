# Časované segmenty pro audio slovíčka

## Účel workflow

Workflow slouží k lokálnímu vygenerování časovaných segmentů pro stránku `/spanelsko-ceska-slovicka`. Produkční web nemá záviset na `ffmpeg`; `ffmpeg` je pouze vývojový nástroj pro přípravu JSON souboru.

Výstupem je JSON objekt, který lze později vložit do `public/sounds/index.json` k odpovídajícímu audio souboru.

## Očekávaná struktura audia

Každá položka v audio souboru má pořadí:

1. španělské slovíčko nebo fráze
2. krátká pauza
3. český překlad
4. delší pauza mezi dvojicemi

Jedna slovíčková položka tedy odpovídá dvěma řečovým blokům:

1. španělsky
2. česky

## Příprava TXT souborů

Zdrojové texty jsou uložené v `public/sounds/source`.

Každý neprázdný řádek je jedna položka. Český a španělský soubor musí mít stejný počet řádků a stejné pořadí položek.

Příklad pro první lekci:

- `public/sounds/source/01_nakupovani_a_jidlo.cs.txt`
- `public/sounds/source/01_nakupovani_a_jidlo.es.txt`

## Příklad spuštění

```bash
python3 scripts/generate_sound_segments.py \
  --audio public/sounds/01_nakupovani_a_jidlo.mp3 \
  --cs public/sounds/source/01_nakupovani_a_jidlo.cs.txt \
  --es public/sounds/source/01_nakupovani_a_jidlo.es.txt \
  --file 01_nakupovani_a_jidlo.mp3 \
  --title "Nakupování a jídlo" \
  --lesson "Lekce 1" \
  --output public/sounds/generated/01_nakupovani_a_jidlo.segments.json
```

## Parametry detekce ticha

`--noise` určuje práh hlasitosti, který `ffmpeg silencedetect` považuje za ticho. Výchozí hodnota je `-35dB`.

`--min-silence` určuje minimální délku ticha v sekundách. Výchozí hodnota je `0.35`.

Když detekce najde příliš mnoho nebo málo řečových bloků, nejčastěji pomůže jemně upravit právě tyto dvě hodnoty.

## Když počet segmentů nesedí

Skript při chybě vypíše:

- počet českých položek
- počet španělských položek
- počet nalezených začátků ticha
- počet nalezených konců ticha
- počet řečových bloků
- počet výsledných segmentů

Pro běžnou lekci, kde jedna položka obsahuje španělskou a českou část, musí platit:

- počet řečových bloků = počet položek krát 2
- počet výsledných segmentů = počet položek

Pokud hodnoty nesedí, zkus upravit `--noise` nebo `--min-silence` a ověřit, zda audio neobsahuje nečekané pauzy nebo ruchy.

## Vložení do index.json

Vygenerovaný soubor `public/sounds/generated/01_nakupovani_a_jidlo.segments.json` obsahuje objekt s polem `segments`.

Později lze tento objekt vložit nebo sloučit do `public/sounds/index.json` k položce:

```json
{
  "file": "01_nakupovani_a_jidlo.mp3",
  "title": "Nakupování a jídlo",
  "lesson": "Lekce 1",
  "segments": []
}
```

Produkční přehrávač pak může časované segmenty používat bez toho, aby na serveru běžel `ffmpeg`.
