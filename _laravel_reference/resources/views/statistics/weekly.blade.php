<x-app-layout>
    <div class="mx-auto w-full sm:w-3/4 p-3 sm:p-4 md:p-6 lg:p-8">

        @php
            $currentLocale = app()->getLocale();
        @endphp

        <!-- Weekly Statistics Header -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6 md:mb-8">
            <h1 class="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold text-gray-800 mb-1 sm:mb-2">{{ __('messages.weekly_statistics_title') }}</h1>
            <p class="text-gray-600 text-[0.65rem] sm:text-xs md:text-sm lg:text-base">{{ __('messages.weekly_statistics_subtitle') }}</p>
        </div>

        <!-- Controls -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
             <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Year Picker for Weekly Stats -->
                <div>
                    <label for="selectedYear" class="block text-xs sm:text-sm font-medium text-gray-700 mb-2">{{ __('messages.monthly_date_label') }}</label>
                    <select id="selectedYear"
                            class="w-full px-3 sm:px-4 py-2 sm:py-2.5 text-xs sm:text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
                        @php
                            $currentYear = date('Y');
                        @endphp
                        @for ($y = 2025; $y <= $currentYear; $y++)
                            <option value="{{ $y }}" {{ $currentYear == $y ? 'selected' : '' }}>{{ $y }}</option>
                        @endfor
                    </select>
                </div>
            </div>
        </div>

        <!-- Weekly Temperature Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.weekly_temperature_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="weeklyTemperatureChart"></canvas>
            </div>
        </div>

        <!-- Weekly Pressure Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.weekly_pressure_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="weeklyPressureChart"></canvas>
            </div>
        </div>

        <!-- Weekly Humidity Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.weekly_humidity_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="weeklyHumidityChart"></canvas>
            </div>
        </div>

        <!-- Weekly Sea Temperature Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.weekly_sea_temperature_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="weeklySeaTemperatureChart"></canvas>
            </div>
        </div>

    </div>

    @vite(['resources/js/statistics-charts.js'])
</x-app-layout>
