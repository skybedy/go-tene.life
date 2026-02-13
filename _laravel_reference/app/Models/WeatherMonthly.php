<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class WeatherMonthly extends Model
{
    protected $table = 'weather_monthly';

    protected $fillable = [
        'year',
        'month',
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
