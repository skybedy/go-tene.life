<?php

namespace App\Http\Controllers;

use App\Models\WeatherHourly;
use App\Models\WeatherDaily;
use Illuminate\Http\Request;
use Carbon\Carbon;

class WeatherController extends Controller
{
    /**
     * Display the weather charts page
     */
    public function charts()
    {
        return view('weather.charts');
    }

    /**
     * Get hourly weather data for charts (from midnight today)
     */
    public function getHourlyData(Request $request)
    {
        // Get date parameter or use today
        $date = $request->input('date', Carbon::today()->format('Y-m-d'));

        $data = WeatherHourly::where('date', $date)
            ->orderBy('hour', 'asc')
            ->get();

        // Format data for Chart.js
        $labels = [];
        $temperatures = [];
        $pressures = [];
        $humidities = [];

        foreach ($data as $record) {
            $labels[] = sprintf('%02d:00', $record->hour);
            $temperatures[] = round((float) $record->avg_temperature, 1);
            $pressures[] = round((float) $record->avg_pressure, 1);
            $humidities[] = round((float) $record->avg_humidity, 1);
        }

        return response()->json([
            'labels' => $labels,
            'datasets' => [
                'temperature' => $temperatures,
                'pressure' => $pressures,
                'humidity' => $humidities,
            ],
        ]);
    }

    /**
     * Get available date range for data
     */
    public function getDateRange()
    {
        $firstRecord = WeatherHourly::orderBy('date', 'asc')->first();
        $lastRecord = WeatherHourly::orderBy('date', 'desc')->first();

        return response()->json([
            'first_date' => $firstRecord ? $firstRecord->date->format('Y-m-d') : null,
            'last_date' => $lastRecord ? $lastRecord->date->format('Y-m-d') : null,
        ]);
    }

    /**
     * Get daily statistics data for charts
     */
    public function getDailyData(Request $request)
    {
        $days = $request->input('days', 7); // Default 7 days
        $endDate = $request->input('end_date', Carbon::today()->format('Y-m-d'));

        $startDate = Carbon::parse($endDate)->subDays($days - 1);

        $data = WeatherDaily::whereBetween('date', [$startDate, $endDate])
            ->orderBy('date', 'asc')
            ->get();

        $labels = [];
        $avgTemp = [];
        $minTemp = [];
        $maxTemp = [];
        $avgPressure = [];
        $avgHumidity = [];
        $seaTemp = [];

        foreach ($data as $record) {
            $labels[] = $record->date->format('j.n.');
            $avgTemp[] = round((float) $record->avg_temperature, 1);
            $minTemp[] = round((float) $record->min_temperature, 1);
            $maxTemp[] = round((float) $record->max_temperature, 1);
            $avgPressure[] = round((float) $record->avg_pressure, 1);
            $avgHumidity[] = round((float) $record->avg_humidity, 1);
            $seaTemp[] = $record->sea_temperature ? round((float) $record->sea_temperature, 1) : null;
        }

        return response()->json([
            'labels' => $labels,
            'datasets' => [
                'avg_temperature' => $avgTemp,
                'min_temperature' => $minTemp,
                'max_temperature' => $maxTemp,
                'avg_pressure' => $avgPressure,
                'avg_humidity' => $avgHumidity,
                'sea_temperature' => $seaTemp,
            ],
            'statistics' => [
                'temperature' => [
                    'avg' => count($avgTemp) > 0 ? round(array_sum($avgTemp) / count($avgTemp), 1) : 0,
                    'min' => count($minTemp) > 0 ? round(min($minTemp), 1) : 0,
                    'max' => count($maxTemp) > 0 ? round(max($maxTemp), 1) : 0,
                ],
                'pressure' => [
                    'avg' => count($avgPressure) > 0 ? round(array_sum($avgPressure) / count($avgPressure), 1) : 0,
                ],
                'humidity' => [
                    'avg' => count($avgHumidity) > 0 ? round(array_sum($avgHumidity) / count($avgHumidity)) : 0,
                ],
                'sea_temperature' => [
                    'avg' => count(array_filter($seaTemp)) > 0 ? round(array_sum(array_filter($seaTemp)) / count(array_filter($seaTemp)), 1) : null,
                ],
            ],
        ]);
    }

    /**
     * Store or update sea temperature for a specific date
     */
    public function storeSeaTemperature(Request $request)
    {
        $validated = $request->validate([
            'date' => 'required|date_format:Y-m-d',
            'temperature' => 'required|numeric|min:-10|max:50',
        ]);

        $date = $validated['date'];
        $temperature = round($validated['temperature'], 1);

        // Find or create the daily record
        $daily = WeatherDaily::firstOrNew(['date' => $date]);
        $daily->sea_temperature = $temperature;
        $daily->save();

        return response()->json([
            'success' => true,
            'message' => 'Sea temperature saved successfully',
            'data' => [
                'date' => $date,
                'sea_temperature' => $temperature,
            ],
        ]);
    }

    /**
     * Get sea temperature for a specific date
     */
    public function getSeaTemperature(Request $request)
    {
        $date = $request->input('date', Carbon::today()->format('Y-m-d'));

        $daily = WeatherDaily::where('date', $date)->first();

        return response()->json([
            'date' => $date,
            'sea_temperature' => $daily ? $daily->sea_temperature : null,
        ]);
    }

    /**
     * Get monthly daily statistics (each day of the month)
     */
    public function getMonthlyDailyData(Request $request)
    {
        $year = $request->input('year', Carbon::now()->year);
        $month = $request->input('month', Carbon::now()->month);

        // Only allow data from November 2025 onwards
        $requestedDate = Carbon::create($year, $month, 1);
        $minDate = Carbon::create(2025, 11, 1);

        if ($requestedDate->lt($minDate)) {
            return response()->json([
                'labels' => [],
                'datasets' => [
                    'avg_temperature' => [],
                    'avg_pressure' => [],
                    'avg_humidity' => [],
                    'sea_temperature' => [],
                ],
                'message' => 'Data available from November 2025 onwards',
            ]);
        }

        $startDate = Carbon::create($year, $month, 1);
        $endDate = $startDate->copy()->endOfMonth();

        $data = WeatherDaily::whereBetween('date', [$startDate, $endDate])
            ->orderBy('date', 'asc')
            ->get();

        $labels = [];
        $avgTemp = [];
        $avgPressure = [];
        $avgHumidity = [];
        $seaTemp = [];

        // Create array with all days of the month
        $daysInMonth = $endDate->day;
        for ($day = 1; $day <= $daysInMonth; $day++) {
            $labels[] = $day;

            // Find data for this day
            $dayData = $data->firstWhere('date', $startDate->copy()->day($day));

            if ($dayData) {
                $avgTemp[] = $dayData->avg_temperature ? round((float) $dayData->avg_temperature, 1) : null;
                $avgPressure[] = $dayData->avg_pressure ? round((float) $dayData->avg_pressure, 1) : null;
                $avgHumidity[] = $dayData->avg_humidity ? round((float) $dayData->avg_humidity, 1) : null;
                $seaTemp[] = $dayData->sea_temperature ? round((float) $dayData->sea_temperature, 1) : null;
            } else {
                $avgTemp[] = null;
                $avgPressure[] = null;
                $avgHumidity[] = null;
                $seaTemp[] = null;
            }
        }

        return response()->json([
            'labels' => $labels,
            'datasets' => [
                'avg_temperature' => $avgTemp,
                'avg_pressure' => $avgPressure,
                'avg_humidity' => $avgHumidity,
                'sea_temperature' => $seaTemp,
            ],
        ]);
    }

    /**
     * Get annual statistics (monthly averages)
     */
    public function getAnnualData(Request $request)
    {
        $year = $request->input('year', Carbon::now()->year);

        // Get all monthly data for the selected year
        $data = \App\Models\WeatherMonthly::where('year', $year)
            ->orderBy('month', 'asc')
            ->get();

        $labels = [];
        $avgTemp = [];
        $avgPressure = [];
        $avgHumidity = [];
        // Sea temperature is not in WeatherMonthly yet, need to aggregate from WeatherDaily or add to table
        // For now, let's aggregate from WeatherDaily for each month
        $seaTemp = [];

        foreach ($data as $record) {
            // Label: Month Year (e.g., 1/2026)
            $labels[] = $record->month . '/' . $record->year;
            
            $avgTemp[] = $record->avg_temperature ? round((float) $record->avg_temperature, 1) : null;
            $avgPressure[] = $record->avg_pressure ? round((float) $record->avg_pressure, 1) : null;
            $avgHumidity[] = $record->avg_humidity ? round((float) $record->avg_humidity, 1) : null;

            // Calculate average sea temperature for this month
            $monthStart = Carbon::create($record->year, $record->month, 1);
            $monthEnd = $monthStart->copy()->endOfMonth();
            
            $avgSeaTemp = WeatherDaily::whereBetween('date', [$monthStart, $monthEnd])
                ->whereNotNull('sea_temperature')
                ->avg('sea_temperature');
            
            $seaTemp[] = $avgSeaTemp ? round($avgSeaTemp, 1) : null;
        }

        return response()->json([
            'labels' => $labels,
            'datasets' => [
                'avg_temperature' => $avgTemp,
                'avg_pressure' => $avgPressure,
                'avg_humidity' => $avgHumidity,
                'sea_temperature' => $seaTemp,
            ],
        ]);
    }

    /**
     * Get weekly statistics
     */
    public function getWeeklyData(Request $request)
    {
        $year = $request->input('year', Carbon::now()->year);

        // Get all weekly data from selected year
        $data = \App\Models\WeatherWeekly::where('year', $year)
            ->orderBy('week', 'asc')
            ->get();

        $labels = [];
        $avgTemp = [];
        $avgPressure = [];
        $avgHumidity = [];
        // Sea temperature aggregation for weeks
        $seaTemp = [];

        foreach ($data as $record) {
            // Calculate week start and end dates
            if ($record->week_start && $record->week_end) {
                $weekStart = $record->week_start;
                $weekEnd = $record->week_end;
            } else {
                 $dto = new Carbon();
                 $dto->setISODate($record->year, $record->week);
                 $weekStart = $dto->startOfWeek();
                 $weekEnd = $dto->copy()->endOfWeek();
            }

            // Label: Week/Year (DateRange) e.g., "42/2025 (13.10. - 19.10.)"
            $labels[] = $record->week . '/' . $record->year . ' (' . $weekStart->format('j.n.') . ' - ' . $weekEnd->format('j.n.') . ')';
            
            $avgTemp[] = $record->avg_temperature ? round((float) $record->avg_temperature, 1) : null;
            $avgPressure[] = $record->avg_pressure ? round((float) $record->avg_pressure, 1) : null;
            $avgHumidity[] = $record->avg_humidity ? round((float) $record->avg_humidity, 1) : null;

            // Calculate average sea temperature for this week
            $avgSeaTemp = WeatherDaily::whereBetween('date', [$weekStart->format('Y-m-d'), $weekEnd->format('Y-m-d')])
                ->whereNotNull('sea_temperature')
                ->avg('sea_temperature');
            
            $seaTemp[] = $avgSeaTemp ? round($avgSeaTemp, 1) : null;
        }

        return response()->json([
            'labels' => $labels,
            'datasets' => [
                'avg_temperature' => $avgTemp,
                'avg_pressure' => $avgPressure,
                'avg_humidity' => $avgHumidity,
                'sea_temperature' => $seaTemp,
            ],
        ]);
    }
}
