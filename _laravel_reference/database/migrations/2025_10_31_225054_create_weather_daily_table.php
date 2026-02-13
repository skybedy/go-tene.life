<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('weather_daily', function (Blueprint $table) {
            $table->id();
            $table->date('date')->unique()->comment('Date of the measurement');
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
            
            $table->index('date');
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('weather_daily');
    }
};