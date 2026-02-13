#!/bin/bash

# Barvy pro lepší přehlednost
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Spouštím deployment TenerLife...${NC}"

# 1. Stáhnutí nejnovějšího kódu
echo -e "${BLUE}📥 Stahuji změny z GitHubu...${NC}"

# Nejdřív musíme Gitu dovolit ty soubory vidět, aby je mohl aktualizovat
git ls-files -z _laravel_reference | xargs -0 git update-index --no-skip-worktree 2>/dev/null
git checkout _laravel_reference 2>/dev/null

git pull origin main

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Chyba při stahování z Gitu!${NC}"
    exit 1
fi

# TEĎ TEN TRIK: Řekneme Gitu, aby ignoroval, že ty soubory smažeme
echo -e "${BLUE}🧹 Čistím server a schovávám změny pro Git...${NC}"
git ls-files -z _laravel_reference | xargs -0 git update-index --skip-worktree 2>/dev/null
rm -rf _laravel_reference

# 2. Build binárky
echo -e "${BLUE}🏗️ Sestavuji novou binárku...${NC}"
go build -o tenelife-app ./main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build byl úspěšný!${NC}"
    echo -e "${BLUE}💡 Tip: Doporučená složka je ~/apps/tene.life${NC}"
    echo -e "${BLUE}💡 Nyní můžeš aplikaci restartovat:${NC}"
    echo -e "   sudo systemctl restart tenelife"
else
    echo -e "${RED}❌ Build selhal!${NC}"
    exit 1
fi
