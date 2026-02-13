# Deployment na produkční server

## 1. Nahrání souborů
Nahrajte všechny soubory projektu na server pomocí FTP/SFTP nebo git pull.

## 2. Konfigurace .env souboru
Na serveru upravte soubor `.env` a ujistěte se, že obsahuje správné nastavení lokalizace:

```env
APP_LOCALE=cs
APP_FALLBACK_LOCALE=cs
APP_FAKER_LOCALE=cs_CZ
```

Také zkontrolujte další důležitá nastavení:
- `APP_ENV=production`
- `APP_DEBUG=false`
- Databázové připojení (DB_*)
- `APP_URL` - URL vaší produkční aplikace

## 3. Vyčištění cache
Po nahrání souborů a úpravě `.env` VŽDY spusťte:

```bash
php artisan config:clear
php artisan cache:clear
php artisan view:clear
php artisan route:clear
```

Nebo použijte připravený skript:
```bash
chmod +x clear-cache.sh
./clear-cache.sh
```

## 4. Optimalizace pro produkci
Pro lepší výkon na produkci:

```bash
php artisan config:cache
php artisan route:cache
php artisan view:cache
```

## 5. Build assets (pokud používáte Vite)
```bash
npm install
npm run build
```

## 6. Kontrola
- Zkontrolujte, zda se překlady zobrazují správně
- Otestujte přepínání jazyků
- Ověřte funkčnost všech stránek (/, /statistics, apod.)

## Časté problémy

### Překlady se nezobrazují (zobrazuje se "messages.home" místo "Webkamera")
**Řešení:**
1. Zkontrolujte `.env` - musí obsahovat `APP_LOCALE=cs`
2. Vyčistěte cache: `php artisan config:clear && php artisan cache:clear`
3. Znovu vytvořte cache: `php artisan config:cache`

### 404 chyba na /statistics
**Řešení:**
1. Vyčistěte route cache: `php artisan route:clear`
2. Zkontrolujte, zda je soubor `routes/web.php` správně nahrán

### Grafy se nezobrazují
**Řešení:**
1. Spusťte `npm run build` pro build assets
2. Zkontrolujte, zda je složka `public/build` nahrána na server
