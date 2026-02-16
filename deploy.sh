#!/bin/bash

# Barvy pro lepší přehlednost
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Spouštím deployment TenerLife...${NC}"

# 1. Stáhnutí nejnovějšího kódu
echo -e "${BLUE}📥 Stahuji změny z GitHubu...${NC}"

# Seznam souborů pro pozdější úklid (tady definujeme, co se má po buildu smazat)
FILES_TO_HIDE="_laravel_reference .air.toml .env.example main.go go.mod go.sum internal views public/js public/build public/storage"

# AGRESIVNÍ OBNOVA: Najdi všechny skryté soubory (skip-worktree) a odhal je
git ls-files -v | grep '^S' | awk '{print $2}' | tr '\n' '\0' | xargs -0 -r git update-index --no-skip-worktree

# Vynutí shodu s repozitářem (obnoví smazané soubory)
git reset --hard HEAD
git pull origin main

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Chyba při stahování z Gitu!${NC}"
    exit 1
fi

# 2. Build frontend a binárky
echo -e "${BLUE}🏗️ Sestavuji frontend a novou binárku...${NC}"

# Instalace závislostí a build Tailwindu
npm install && npm run build
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Chyba při buildu Tailwindu!${NC}"
    exit 1
fi

go mod tidy
go build -o tenelife-app .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build byl úspěšný!${NC}"
    
    # TEĎ TEN TRIK: Úklid všeho nepotřebného po úspěšném buildu
    echo -e "${BLUE}🧹 Čistím server od zdrojových kódů (Production mode)...${NC}"
    
    # Přidány i nové soubory Tailwindu k úklidu
    FILES_TO_HIDE="$FILES_TO_HIDE resources package.json package-lock.json node_modules public/css"
    
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
