<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\IndexController;
use App\Http\Controllers\WeatherController;
use App\Http\Controllers\CameraUploadController;

// API endpoints for weather data (no locale prefix)
Route::get('/api/weather/hourly', [WeatherController::class, 'getHourlyData'])->name('api.weather.hourly');
Route::get('/api/weather/daily', [WeatherController::class, 'getDailyData'])->name('api.weather.daily');
Route::get('/api/weather/monthly-daily', [WeatherController::class, 'getMonthlyDailyData'])->name('api.weather.monthly-daily');
Route::get('/api/weather/weekly', [WeatherController::class, 'getWeeklyData'])->name('api.weather.weekly');
Route::get('/api/weather/annual', [WeatherController::class, 'getAnnualData'])->name('api.weather.annual');
Route::get('/api/weather/sea-temperature', [WeatherController::class, 'getSeaTemperature'])->name('api.weather.sea-temperature.get');
Route::post('/api/weather/sea-temperature', [WeatherController::class, 'storeSeaTemperature'])->name('api.weather.sea-temperature.store');

// Routes with locale prefix
Route::group(['prefix' => '{locale}', 'where' => ['locale' => 'cs|en|es|de|it|pl|hu|fr']], function () {
    Route::get('/', [IndexController::class, 'index'])->name('index.index.locale');
    Route::get('/webcam/big', [IndexController::class, 'webcamBig'])->name('index.webcam.big.locale');
    Route::get('/statistics', [IndexController::class, 'statistics'])->name('index.statistics.locale');
    Route::get('/statistics/daily', [IndexController::class, 'dailyStatistics'])->name('index.statistics.daily.locale');
    Route::get('/statistics/monthly', [IndexController::class, 'monthlyStatistics'])->name('index.statistics.monthly.locale');
    Route::get('/statistics/weekly', [IndexController::class, 'weeklyStatistics'])->name('index.statistics.weekly.locale');
    Route::get('/statistics/annual', [IndexController::class, 'annualStatistics'])->name('index.statistics.annual.locale');
});

// Default routes (Czech - no prefix)
Route::get('/', [IndexController::class, 'index'])->name('index.index');
Route::get('/webcam/big', [IndexController::class, 'webcamBig'])->name('index.webcam.big');
Route::get('/statistics', [IndexController::class, 'statistics'])->name('index.statistics'); // Keep for backward compatibility or redirect
Route::get('/statistics/daily', [IndexController::class, 'dailyStatistics'])->name('index.statistics.daily');
Route::get('/statistics/monthly', [IndexController::class, 'monthlyStatistics'])->name('index.statistics.monthly');
Route::get('/statistics/weekly', [IndexController::class, 'weeklyStatistics'])->name('index.statistics.weekly');
Route::get('/statistics/annual', [IndexController::class, 'annualStatistics'])->name('index.statistics.annual');
