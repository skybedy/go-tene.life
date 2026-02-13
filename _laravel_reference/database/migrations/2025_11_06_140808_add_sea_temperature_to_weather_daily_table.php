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
            $table->decimal('sea_temperature', 5, 1)->nullable()->after('date')->comment('Sea water temperature in °C (manually measured)');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('weather_daily', function (Blueprint $table) {
            $table->dropColumn('sea_temperature');
        });
    }
};
