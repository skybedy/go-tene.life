#!/bin/bash
# Script to clear Laravel cache on production server

echo "Clearing Laravel cache..."

# Clear all caches
php artisan config:clear
php artisan cache:clear
php artisan route:clear
php artisan view:clear

# Rebuild cache for production
php artisan config:cache
php artisan route:cache
php artisan view:cache

echo "Cache cleared and rebuilt successfully!"
