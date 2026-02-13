#!/bin/bash

# Barvy pro lepší přehlednost
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Spouštím deployment TenerLife...${NC}"

# 1. Stáhnutí nejnovějšího kódu
echo -e "${BLUE}📥 Stahuji změny z GitHubu...${NC}"

# Seznam souborů a složek, které v produkci nechceme, ale jsou v Gitu
FILES_TO_HIDE="_laravel_reference .air.toml .env.example main.go go.mod go.sum internal views"

# Nejdřív musíme Gitu dovolit ty soubory vidět, aby je mohl aktualizovat
for FILE in $FILES_TO_HIDE; do
    git ls-files -z "$FILE" | xargs -0 git update-index --no-skip-worktree 2>/dev/null
    git checkout "$FILE" 2>/dev/null
done

git pull origin main

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Chyba při stahování z Gitu!${NC}"
    exit 1
fi

# 2. Build binárky (musí proběhnout dokud jsou soubory na disku)
echo -e "${BLUE}🏗️ Sestavuji novou binárku...${NC}"
go build -o tenelife-app ./main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build byl úspěšný!${NC}"
    
    # TEĎ TEN TRIK: Úklid všeho nepotřebného po úspěšném buildu
    echo -e "${BLUE}🧹 Čistím server od zdrojových kódů (Production mode)...${NC}"
    for FILE in $FILES_TO_HIDE; do
        git ls-files -z "$FILE" | xargs -0 git update-index --skip-worktree 2>/dev/null
        rm -rf "$FILE"
    done

    echo -e "${BLUE}💡 Tip: Doporučená složka je ~/apps/tene.life${NC}"
    echo -e "${BLUE}💡 Nyní můžeš aplikaci restartovat:${NC}"
    echo -e "   sudo systemctl restart tenelife"
else
    echo -e "${RED}❌ Build selhal! Zdrojové kódy ponechány pro diagnostiku.${NC}"
    exit 1
fi
