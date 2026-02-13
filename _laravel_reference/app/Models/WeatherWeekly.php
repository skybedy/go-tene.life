<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class WeatherWeekly extends Model
{
    protected $table = 'weather_weekly';

    protected $fillable = [
        'year',
        'week',
        'week_start',
        'week_end',
        'avg_temperature',
        'min_temperature',
        'max_temperature',
        'avg_pressure',
        'min_pressure',
        'max_pressure',
        'avg_humidity',
        'min_humidity',
        'max_humidity',
        'samples_count',
    ];

    protected $casts = [
        'week_start' => 'date',
        'week_end' => 'date',
        'avg_temperature' => 'float',
        'min_temperature' => 'float',
        'max_temperature' => 'float',
        'avg_pressure' => 'float',
        'min_pressure' => 'float',
        'max_pressure' => 'float',
        'avg_humidity' => 'float',
        'min_humidity' => 'float',
        'max_humidity' => 'float',
    ];
}
