<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\WeatherDaily;
use App\Models\WeatherWeekly;
use App\Models\WeatherMonthly;

class IndexController extends Controller
{
    public function index()
    {
        $weatherData = $this->getWeatherData();
        $seaTemperature = $this->getSeaTemperature();

        return view('index',[
            'weatherData' => $weatherData,
            'weatherTimestamp' => $weatherData ? $this->getWeatherTimestamp() : null,
            'seaTemperature' => $seaTemperature
        ]);
    }

    private function getWeatherData()
    {

        $filePath = public_path('files/weather.json');

        if (!file_exists($filePath)) {
            return null; // nebo []
        }

        $rawData = file_get_contents($filePath);

        $data = json_decode($rawData, true);

        return $data;
    }

    private function getWeatherTimestamp()
    {
        $filePath = public_path('files/weather.json');

        if (!file_exists($filePath)) {
            return null;
        }

        return filemtime($filePath);
    }

    private function getSeaTemperature()
    {
        // Try to get today's sea temperature first
        $today = WeatherDaily::whereDate('date', today())
            ->whereNotNull('sea_temperature')
            ->first();

        if ($today && $today->sea_temperature) {
            return $today->sea_temperature;
        }

        // If not found, get the most recent sea temperature
        $latest = WeatherDaily::whereNotNull('sea_temperature')
            ->orderBy('date', 'desc')
            ->first();

        return $latest ? $latest->sea_temperature : null;
    }



    public function webcamBig()
    {
        return view('webcam-big');
    }

    public function statistics()
    {
        return redirect()->route('index.statistics.daily');
    }

    public function dailyStatistics()
    {
        // Get daily statistics (last 30 days)
        $dailyStats = WeatherDaily::orderBy('date', 'desc')
            ->limit(30)
            ->get();

        return view('statistics.daily', [
            'dailyStats' => $dailyStats,
        ]);
    }

    public function monthlyStatistics()
    {
        // Get monthly statistics (last 12 months) for the table
        $monthlyStats = WeatherMonthly::orderBy('year', 'desc')
            ->orderBy('month', 'desc')
            ->limit(12)
            ->get();

        return view('statistics.monthly', [
            'monthlyStats' => $monthlyStats,
        ]);
    }

    public function weeklyStatistics()
    {
        // Get weekly statistics ordered by date for the table
        $weeklyStats = WeatherWeekly::orderBy('year', 'desc')
            ->orderBy('week', 'desc')
            ->get();

        return view('statistics.weekly', [
            'weeklyStats' => $weeklyStats,
        ]);
    }

    public function annualStatistics()
    {
        // Get all monthly statistics ordered by date for the table
        $annualStats = WeatherMonthly::orderBy('year', 'desc')
            ->orderBy('month', 'desc')
            ->get();

        return view('statistics.annual', [
            'annualStats' => $annualStats,
        ]);
    }
}
