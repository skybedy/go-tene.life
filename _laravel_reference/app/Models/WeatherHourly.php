<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class WeatherHourly extends Model
{
    /**
     * The table associated with the model.
     *
     * @var string
     */
    protected $table = 'weather_hourly';

    /**
     * The attributes that are mass assignable.
     *
     * @var array<int, string>
     */
    protected $fillable = [
        'date',
        'hour',
        'avg_temperature',
        'avg_pressure',
        'avg_humidity',
        'samples_count',
    ];

    /**
     * The attributes that should be cast.
     *
     * @var array<string, string>
     */
    protected $casts = [
        'date' => 'date',
        'hour' => 'integer',
        'avg_temperature' => 'decimal:1',
        'avg_pressure' => 'decimal:1',
        'avg_humidity' => 'decimal:1',
        'samples_count' => 'integer',
    ];

    /**
     * Scope to filter by date range
     */
    public function scopeDateRange($query, $start, $end)
    {
        return $query->whereBetween('date', [$start, $end]);
    }

    /**
     * Scope to get data for specific date
     */
    public function scopeForDate($query, $date)
    {
        return $query->where('date', $date);
    }
}
