#!/bin/bash

# Barvy pro lepší přehlednost
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Spouštím deployment TenerLife...${NC}"

# 1. Stáhnutí nejnovějšího kódu
echo -e "${BLUE}📥 Stahuji změny z GitHubu...${NC}"
git pull origin main

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Chyba při stahování z Gitu!${NC}"
    exit 1
fi

# 2. Build binárky
echo -e "${BLUE}🏗️ Sestavuji novou binárku...${NC}"
go build -o tenelife-app ./main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build byl úspěšný!${NC}"
    echo -e "${BLUE}💡 Nyní můžeš aplikaci spustit nebo restartovat službu.${NC}"
    echo -e "   Příklad: ./tenelife-app"
else
    echo -e "${RED}❌ Build selhal!${NC}"
    exit 1
fi
