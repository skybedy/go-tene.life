<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('weather_weekly', function (Blueprint $table) {
            $table->id();
            $table->integer('year')->comment('Year (e.g., 2025)');
            $table->tinyInteger('week')->comment('ISO week number (1-53)');
            $table->date('week_start')->comment('Monday of the week');
            $table->date('week_end')->comment('Sunday of the week');
            $table->decimal('avg_temperature', 5, 1)->comment('Average temperature in °C');
            $table->decimal('min_temperature', 5, 1)->comment('Minimum temperature in °C');
            $table->decimal('max_temperature', 5, 1)->comment('Maximum temperature in °C');
            $table->decimal('avg_pressure', 7, 1)->comment('Average pressure in hPa');
            $table->decimal('min_pressure', 7, 1)->comment('Minimum pressure in hPa');
            $table->decimal('max_pressure', 7, 1)->comment('Maximum pressure in hPa');
            $table->decimal('avg_humidity', 5, 1)->comment('Average humidity in %');
            $table->decimal('min_humidity', 5, 1)->comment('Minimum humidity in %');
            $table->decimal('max_humidity', 5, 1)->comment('Maximum humidity in %');
            $table->integer('samples_count')->comment('Number of measurements used');
            $table->timestamp('created_at')->default(DB::raw('CURRENT_TIMESTAMP'));
            $table->timestamp('updated_at')->default(DB::raw('CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP'));
            
            $table->unique(['year', 'week']);
            $table->index(['year', 'week']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('weather_weekly');
    }
};