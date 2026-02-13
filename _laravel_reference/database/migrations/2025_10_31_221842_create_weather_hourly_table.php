<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('weather_hourly', function (Blueprint $table) {
            $table->id();
            $table->date('date')->comment('Date of the measurement');
            $table->tinyInteger('hour')->comment('Hour (0-23)');
            $table->decimal('avg_temperature', 5, 1)->comment('Average temperature in °C');
            $table->decimal('avg_pressure', 7, 1)->comment('Average pressure in hPa');
            $table->decimal('avg_humidity', 5, 1)->comment('Average humidity in %');
            $table->integer('samples_count')->comment('Number of measurements used for average');
            $table->timestamp('created_at')->default(DB::raw('CURRENT_TIMESTAMP'));
            $table->timestamp('updated_at')->default(DB::raw('CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP'));
            
            $table->unique(['date', 'hour']);
            $table->index('date');
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('weather_hourly');
    }
};