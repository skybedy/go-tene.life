<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::table('weather_daily', function (Blueprint $table) {
            $table->decimal('avg_temperature', 5, 1)->nullable()->change();
            $table->decimal('min_temperature', 5, 1)->nullable()->change();
            $table->decimal('max_temperature', 5, 1)->nullable()->change();
            $table->decimal('avg_pressure', 7, 1)->nullable()->change();
            $table->decimal('min_pressure', 7, 1)->nullable()->change();
            $table->decimal('max_pressure', 7, 1)->nullable()->change();
            $table->decimal('avg_humidity', 5, 1)->nullable()->change();
            $table->decimal('min_humidity', 5, 1)->nullable()->change();
            $table->decimal('max_humidity', 5, 1)->nullable()->change();
            $table->integer('samples_count')->nullable()->change();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('weather_daily', function (Blueprint $table) {
            $table->decimal('avg_temperature', 5, 1)->nullable(false)->change();
            $table->decimal('min_temperature', 5, 1)->nullable(false)->change();
            $table->decimal('max_temperature', 5, 1)->nullable(false)->change();
            $table->decimal('avg_pressure', 7, 1)->nullable(false)->change();
            $table->decimal('min_pressure', 7, 1)->nullable(false)->change();
            $table->decimal('max_pressure', 7, 1)->nullable(false)->change();
            $table->decimal('avg_humidity', 5, 1)->nullable(false)->change();
            $table->decimal('min_humidity', 5, 1)->nullable(false)->change();
            $table->decimal('max_humidity', 5, 1)->nullable(false)->change();
            $table->integer('samples_count')->nullable(false)->change();
        });
    }
};
