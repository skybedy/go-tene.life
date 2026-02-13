# AXIS Camera Upload - Stav implementace

**Datum:** 9.11.2025 23:50
**Kamera:** AXIS P1355 (firmware 2020)
**Server:** tene.life (Ubuntu 22.04, OpenSSL 3.0.2)

## ✅ Co funguje:

1. **HTTP endpoint přijímá obrázky**
   - URL: `http://tene.life:8080/api/camera/upload`
   - Metody: GET, POST, PUT
   - Ukládá jako: `/var/www/laravel-tene.life/public/images/tenelife.jpg`
   - Vždy přepíše předchozí soubor (bez timestampu)
   - Testováno curl - funguje ✅

2. **Nginx listener na portu 8080**
   - Soubor: `/etc/nginx/sites-available/tene.life`
   - Server block: `listen 8080;`
   - PHP-FPM: funguje

3. **Laravel routing**
   - Soubor: `/var/www/laravel-tene.life/routes/api.php`
   - Route: `Route::match(['post', 'put', 'get'], '/camera/upload', ...)`

4. **Firewall**
   - Port 8080: OTEVŘENÝ (ufw allow 8080)
   - Port 8443: OTEVŘENÝ (pro stunnel, ale nepoužívá se)

5. **File permissions**
   - Složka images: vlastník `www-data:www-data`
   - PHP může zapisovat ✅

## ⚠️ Co NEFUNGUJE / CHYBÍ:

### 1. **AUTENTIZACE - PRIORITA #1**
   - ❌ Endpoint je BEZ autentizace - kdokoliv může uploadovat!
   - V kódu je autentizace zakomentovaná (řádky 27-41 v CameraUploadController.php)
   - **Problém:** HTTP Basic Auth nefunguje - kamera v HTTP Recipient módu credentials neodesílá

   **Řešení pro zítřek:**
   - Implementovat autentizaci pomocí tokenu v URL
   - Příklad: `http://tene.life:8080/api/camera/upload?token=TAJNY_TOKEN`
   - Token uložit do `.env` jako `CAMERA_UPLOAD_TOKEN`
   - Otestovat s kamerou

### 2. **Kamera neposílá obrázky**
   - Při testu kamera hlásí "Upload successful"
   - Ale obrázek se neukládá
   - **Možné příčiny:**
     - Kamera posílá prázdné tělo (testovací režim?)
     - Action Rule není správně nakonfigurované
     - Kamera čeká na trigger událost (pohyb, schedule, atd.)

   **K ověření zítřek:**
   - Zkontrolovat Action Rule v kameře
   - Ověřit, že kamera posílá skutečná data (ne jen test)
   - Zkontrolovat debug log: `/tmp/camera_debug.log`

### 3. **Debug kód v produkci**
   - Řádky 15-25 v CameraUploadController.php
   - Loguje každý požadavek do `/tmp/camera_debug.log`
   - **TODO:** Po otestování odstranit nebo přesunout za `if (env('APP_DEBUG'))`

## 📋 TODO pro zítřek:

### Priorita 1: Zabezpečení
- [ ] Přidat autentizaci pomocí tokenu v URL
- [ ] Vygenerovat silný token (např. 64 znaků hex)
- [ ] Přidat do `.env`: `CAMERA_UPLOAD_TOKEN=...`
- [ ] Upravit controller - zkontrolovat token z `$request->query('token')`
- [ ] Nastavit URL v kameře: `http://tene.life:8080/api/camera/upload?token=TOKEN`

### Priorita 2: Otestovat s kamerou
- [ ] Zkontrolovat Action Rule v kameře
- [ ] Spustit test upload z kamery
- [ ] Zkontrolovat `/tmp/camera_debug.log` - co kamera posílá?
- [ ] Ověřit, že se vytváří nový `tenelife.jpg`

### Priorita 3: Cleanup
- [ ] Odstranit debug kód nebo chránit za APP_DEBUG
- [ ] Rozhodnout: zůstat na portu 8080 nebo přejít na 80?
- [ ] Zavřít port 8443 pokud stunnel nepoužíváme
- [ ] Vypnout stunnel service pokud nepoužíváme

### Volitelné:
- [ ] Zvážit návrat k FTP (bylo bezpečnější s credentials)
- [ ] Přidat rate limiting na endpoint
- [ ] Přidat monitoring - notifikace když upload selže

## 🔧 Důležité soubory:

```
/var/www/laravel-tene.life/app/Http/Controllers/CameraUploadController.php
/var/www/laravel-tene.life/routes/api.php
/var/www/laravel-tene.life/.env
/etc/nginx/sites-available/tene.life
/etc/stunnel/axis-camera.conf
/tmp/camera_debug.log (pokud existuje)
```

## 📝 Poznámky:

- **HTTPS nefunguje:** OpenSSL 3.0.2 nepodporuje SSLv3 které kamera používá
- **Stunnel nefunguje:** Také používá OpenSSL 3.0, stejný problém
- **HTTP Basic Auth nefunguje:** Kamera v HTTP Recipient módu credentials neposílá
- **Port 8080:** Používáme protože tam je Nginx listener, ale není nutné - můžeme přejít na port 80

## ❓ Otázky k rozmyšlení:

1. Je HTTP endpoint s tokenem bezpečnější než FTP s username+password?
2. Chceme zůstat na portu 8080 nebo přejít na standardní port 80?
3. Chceme ukládat historii obrázků nebo vždy jen poslední?

---

**Status:** Endpoint funguje technicky, ale bez autentizace = NEBEZPEČNÉ pro produkci!
**Next step:** Přidat token autentizaci a otestovat s kamerou.
