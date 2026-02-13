<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class WeatherDaily extends Model
{
    protected $table = 'weather_daily';

    protected $fillable = [
        'date',
        'sea_temperature',
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
        'date' => 'date',
        'sea_temperature' => 'float',
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
