<x-app-layout>
    <div class="mx-auto w-full sm:w-3/4 p-3 sm:p-4 md:p-6 lg:p-8">

        @php
            $currentLocale = app()->getLocale();
        @endphp

        <!-- Page Header -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6 md:mb-8">
            <h1 class="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold text-gray-800 mb-1 sm:mb-2">{{ __('messages.statistics_title') }}</h1>
            <p class="text-gray-600 text-[0.65rem] sm:text-xs md:text-sm lg:text-base">{{ __('messages.statistics_subtitle') }}</p>
        </div>

        <!-- Controls -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Date Picker for Hourly Stats -->
                <div>
                    <label for="selectedDate" class="block text-xs sm:text-sm font-medium text-gray-700 mb-2">{{ __('messages.hourly_date_label') }}</label>
                    <input type="date"
                           id="selectedDate"
                           class="w-full px-3 sm:px-4 py-2 sm:py-2.5 text-xs sm:text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                           value="{{ date('Y-m-d') }}"
                           min="2025-11-01"
                           max="{{ date('Y-m-d') }}">
                </div>

                <!-- Month Picker for Monthly Stats -->
                <div>
                    <label for="selectedMonth" class="block text-xs sm:text-sm font-medium text-gray-700 mb-2">{{ __('messages.monthly_date_label') }}</label>
                    <div class="grid grid-cols-2 gap-2">
                        <select id="selectedMonthNum"
                                class="w-full px-3 sm:px-4 py-2 sm:py-2.5 text-xs sm:text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
                            @for ($m = 1; $m <= 12; $m++)
                                <option value="{{ $m }}" {{ date('n') == $m ? 'selected' : '' }}>
                                    {{ __('messages.month_' . $m) }}
                                </option>
                            @endfor
                        </select>
                        <select id="selectedYear"
                                class="w-full px-3 sm:px-4 py-2 sm:py-2.5 text-xs sm:text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
                            @php
                                $currentYear = date('Y');
                                $currentMonth = date('n');
                            @endphp
                            @for ($y = 2025; $y <= $currentYear; $y++)
                                <option value="{{ $y }}" {{ $currentYear == $y ? 'selected' : '' }}>{{ $y }}</option>
                            @endfor
                        </select>
                    </div>
                </div>
            </div>
        </div>

        <!-- Statistics Summary -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 sm:gap-4 mb-4 sm:mb-6">
            <!-- Temperature Stats -->
            <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4">
                <h3 class="text-xs sm:text-sm font-semibold text-gray-700 mb-2">{{ __('messages.temperature_chart') }}</h3>
                <div class="space-y-1 text-[0.65rem] sm:text-xs">
                    <div class="flex justify-between">
                        <span class="text-gray-600">{{ __('messages.average') }}</span>
                        <span class="font-medium" id="stat-temp-avg">--</span>
                    </div>
                    <div class="flex justify-between">
                        <span class="text-gray-600">{{ __('messages.min') }}</span>
                        <span class="font-medium" id="stat-temp-min">--</span>
                    </div>
                    <div class="flex justify-between">
                        <span class="text-gray-600">{{ __('messages.max') }}</span>
                        <span class="font-medium" id="stat-temp-max">--</span>
                    </div>
                </div>
            </div>

            <!-- Pressure Stats -->
            <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4">
                <h3 class="text-xs sm:text-sm font-semibold text-gray-700 mb-2">{{ __('messages.pressure_chart') }}</h3>
                <div class="space-y-1 text-[0.65rem] sm:text-xs">
                    <div class="flex justify-between">
                        <span class="text-gray-600">{{ __('messages.average') }}</span>
                        <span class="font-medium" id="stat-pressure-avg">--</span>
                    </div>
                    <div class="flex justify-between">
                        <span class="text-gray-600">{{ __('messages.min') }}</span>
                        <span class="font-medium" id="stat-pressure-min">--</span>
                    </div>
                    <div class="flex justify-between">
                        <span class="text-gray-600">{{ __('messages.max') }}</span>
                        <span class="font-medium" id="stat-pressure-max">--</span>
                    </div>
                </div>
            </div>

            <!-- Humidity Stats -->
            <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4">
                <h3 class="text-xs sm:text-sm font-semibold text-gray-700 mb-2">{{ __('messages.humidity_chart') }}</h3>
                <div class="space-y-1 text-[0.65rem] sm:text-xs">
                    <div class="flex justify-between">
                        <span class="text-gray-600">{{ __('messages.average') }}</span>
                        <span class="font-medium" id="stat-humidity-avg">--</span>
                    </div>
                    <div class="flex justify-between">
                        <span class="text-gray-600">{{ __('messages.min') }}</span>
                        <span class="font-medium" id="stat-humidity-min">--</span>
                    </div>
                    <div class="flex justify-between">
                        <span class="text-gray-600">{{ __('messages.max') }}</span>
                        <span class="font-medium" id="stat-humidity-max">--</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Temperature Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.temperature_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="temperatureChart"></canvas>
            </div>
        </div>

        <!-- Pressure Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.pressure_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="pressureChart"></canvas>
            </div>
        </div>

        <!-- Humidity Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.humidity_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="humidityChart"></canvas>
            </div>
        </div>

        <!-- Monthly Statistics Header -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6 mt-8">
            <h2 class="text-base sm:text-lg md:text-xl lg:text-2xl font-bold text-gray-800 mb-1">{{ __('messages.monthly_statistics_title') }}</h2>
            <p class="text-gray-600 text-[0.65rem] sm:text-xs md:text-sm">{{ __('messages.monthly_statistics_subtitle') }}</p>
        </div>

        <!-- Monthly Temperature Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.monthly_temperature_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="monthlyTemperatureChart"></canvas>
            </div>
        </div>

        <!-- Monthly Pressure Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.monthly_pressure_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="monthlyPressureChart"></canvas>
            </div>
        </div>

        <!-- Monthly Humidity Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.monthly_humidity_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="monthlyHumidityChart"></canvas>
            </div>
        </div>

        <!-- Monthly Sea Temperature Chart -->
        <div class="bg-white/30 rounded-2xl shadow-lg backdrop-blur-sm border border-gray-300 p-3 sm:p-4 md:p-5 lg:p-6 mb-4 sm:mb-6">
            <h2 class="text-sm sm:text-base md:text-lg font-semibold mb-3 sm:mb-4 text-gray-700">{{ __('messages.monthly_sea_temperature_chart') }}</h2>
            <div class="relative h-40 sm:h-48 md:h-56 lg:h-64">
                <canvas id="monthlySeaTemperatureChart"></canvas>
            </div>
        </div>

    </div>

    @vite(['resources/js/statistics-charts.js'])
</x-app-layout>
