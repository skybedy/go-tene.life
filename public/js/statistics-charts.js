document.addEventListener('DOMContentLoaded', function () {
    const i18n = window.statsI18n || {};
    const path = window.location.pathname;
    const isDailyLikePage = path.includes('/statistics/daily') || path.includes('/statistics/recent');
    const charts = {
        temperature: null,
        seaTemperature: null,
        pressure: null,
        humidity: null
    };

    function initChart(id, label, color, unit, isMulti = false) {
        const el = document.getElementById(id);
        if (!el) return null;
        
        const ctx = el.getContext('2d');
        
        const datasets = [];
        if (isMulti) {
            datasets.push(
                {
                    label: i18n.min || 'Min',
                    data: [],
                    borderColor: '#3b82f6',
                    borderWidth: 1,
                    fill: false,
                    tension: 0.4,
                    pointRadius: 0
                },
                {
                    label: i18n.average || 'Average',
                    data: [],
                    borderColor: color,
                    backgroundColor: color + '20',
                    borderWidth: 3,
                    fill: true,
                    tension: 0.4,
                    pointRadius: 2
                },
                {
                    label: i18n.max || 'Max',
                    data: [],
                    borderColor: '#ef4444',
                    borderWidth: 1,
                    fill: false,
                    tension: 0.4,
                    pointRadius: 0
                }
            );
        } else {
            datasets.push({
                label: label,
                data: [],
                borderColor: color,
                backgroundColor: color + '20',
                borderWidth: 2,
                fill: true,
                tension: 0.4,
                pointRadius: 2
            });
        }

        return new Chart(ctx, {
            type: 'line',
            data: {
                labels: [],
                datasets: datasets
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        display: true,
                        grid: { display: false }
                    },
                    y: {
                        display: true,
                        ticks: {
                            callback: function(value) {
                                const numeric = Number(value);
                                if (Number.isNaN(numeric)) return value + unit;
                                const rounded = Math.round((numeric + Number.EPSILON) * 10) / 10;
                                return rounded.toFixed(1) + unit;
                            }
                        }
                    }
                },
                plugins: {
                    legend: { display: isMulti },
                    tooltip: {
                        mode: 'index',
                        intersect: false
                    }
                }
            }
        });
    }

    // Initialize charts if they exist on page
    charts.temperature = initChart('temperatureChart', i18n.temperature || 'Temperature', '#ef4444', '°C', isDailyLikePage);
    charts.seaTemperature = initChart('seaTemperatureChart', i18n.seaTemperature || 'Sea Temperature', '#0ea5e9', '°C');
    charts.pressure = initChart('pressureChart', i18n.pressure || 'Pressure', '#3b82f6', ' hPa');
    charts.humidity = initChart('humidityChart', i18n.humidity || 'Humidity', '#10b981', '%');

    async function loadDailyStats() {
        if (!charts.temperature) return;

        try {
            const days = Number(window.statsDays) || 7;
            const response = await fetch(`/api/weather/daily?days=${encodeURIComponent(days)}`);
            const data = await response.json();

            if (data.labels && data.datasets) {
                // Temperature (Min/Avg/Max)
                charts.temperature.data.labels = data.labels;
                charts.temperature.data.datasets[0].data = data.datasets.min_temperature;
                charts.temperature.data.datasets[1].data = data.datasets.avg_temperature;
                charts.temperature.data.datasets[2].data = data.datasets.max_temperature;
                charts.temperature.update();

                await loadSeaTemperatureHistoryForPath();

                // Pressure
                charts.pressure.data.labels = data.labels;
                charts.pressure.data.datasets[0].data = data.datasets.avg_pressure;
                charts.pressure.update();

                // Humidity
                charts.humidity.data.labels = data.labels;
                charts.humidity.data.datasets[0].data = data.datasets.avg_humidity;
                charts.humidity.update();

                // Update summary boxes
                updateSummary(data.datasets);
            }
        } catch (error) {
            console.error('Error loading daily stats:', error);
        }
    }

    function updateSummary(datasets) {
        if (!datasets) return;

        const avg = arr => {
            if (!Array.isArray(arr)) return null;
            const values = arr.filter(value => typeof value === 'number' && Number.isFinite(value));
            return values.length ? values.reduce((a, b) => a + b, 0) / values.length : null;
        };

        const tempAvg = avg(datasets.avg_temperature);
        const pressAvg = avg(datasets.avg_pressure);
        const humAvg = avg(datasets.avg_humidity);
        const seaTempAvg = avg(datasets.sea_temperature);

        if (document.getElementById('stat-temp-avg')) document.getElementById('stat-temp-avg').textContent = tempAvg === null ? '--' : tempAvg.toFixed(1) + ' °C';
        if (document.getElementById('stat-pressure-avg')) document.getElementById('stat-pressure-avg').textContent = pressAvg === null ? '--' : pressAvg.toFixed(1) + ' hPa';
        if (document.getElementById('stat-humidity-avg')) document.getElementById('stat-humidity-avg').textContent = humAvg === null ? '--' : Math.round(humAvg) + ' %';
        if (document.getElementById('stat-sea-temp-avg')) document.getElementById('stat-sea-temp-avg').textContent = seaTempAvg === null ? '--' : seaTempAvg.toFixed(1) + ' °C';
    }

    function capitalizeFirst(value) {
        if (typeof value !== 'string' || value.length === 0) return value;
        return value.charAt(0).toLocaleUpperCase(document.documentElement.lang || 'cs') + value.slice(1);
    }

    function stripWrappingQuotes(value) {
        if (typeof value !== 'string') return '';
        return value.trim().replace(/^["'“”„]+|["'“”„]+$/g, '');
    }

    function dateFromISODate(dateString) {
        if (typeof dateString !== 'string' || !dateString) return '';
        const parts = dateString.split('-').map(Number);
        if (parts.length !== 3 || parts.some(Number.isNaN)) return null;
        return new Date(Date.UTC(parts[0], parts[1] - 1, parts[2]));
    }

    function formatDateForLocale(dateString) {
        const date = dateFromISODate(dateString);
        if (!date) return dateString || '';

        const locale = document.documentElement.lang || 'cs';
        return new Intl.DateTimeFormat(locale, {
            day: 'numeric',
            month: 'numeric',
            timeZone: 'UTC'
        }).format(date);
    }

    function formatMonthForLocale(dateString) {
        const date = dateFromISODate(dateString);
        if (!date) return '';

        const locale = document.documentElement.lang || 'cs';
        const month = new Intl.DateTimeFormat(locale, {
            month: 'long',
            timeZone: 'UTC'
        }).format(date);
        return capitalizeFirst(month);
    }

    function updateMonthlySummaryTitle(throughDate) {
        const title = document.getElementById('monthlyCurrentSummaryTitle');
        if (!title) return;

        const formattedDate = formatDateForLocale(throughDate);
        const formattedMonth = formatMonthForLocale(throughDate);
        const prefix = stripWrappingQuotes(i18n.monthlyCurrentProgressTo || i18n.monthlyCurrentProgress || '');
        title.textContent = formattedDate && formattedMonth
            ? `${formattedMonth} ${prefix} ${formattedDate}`
            : (i18n.monthlyCurrentProgress || title.textContent);
    }

    // Route handling
    if (path.includes('/statistics/daily')) {
        loadDailyStats();
    } else if (path.includes('/statistics/recent')) {
        loadDailyStats();
    } else if (path.includes('/statistics/weekly')) {
        // Implement weekly...
        loadGenericStats('weekly');
    } else if (path.includes('/statistics/monthly-daily')) {
        loadMonthlyDailyStats();
    } else if (path.includes('/statistics/monthly')) {
        loadGenericStats('monthly', false);
        loadMonthlySummary();
    } else if (path.includes('/statistics/annual')) {
        loadGenericStats('annual');
    }

    async function loadMonthlyDailyStats() {
        if (!charts.temperature) return;

        try {
            const year = window.statsSelectedYear;
            const month = window.statsSelectedMonth;
            const response = await fetch(`/api/weather/monthly-daily?year=${encodeURIComponent(year)}&month=${encodeURIComponent(month)}`);
            const data = await response.json();

            if (data.labels && data.datasets) {
                charts.temperature.data.labels = data.labels;
                charts.temperature.data.datasets[0].data = data.datasets.avg_temperature;
                charts.temperature.update();

                if (charts.seaTemperature) {
                    charts.seaTemperature.data.labels = data.labels;
                    charts.seaTemperature.data.datasets[0].data = data.datasets.sea_temperature;
                    charts.seaTemperature.update();
                }

                charts.pressure.data.labels = data.labels;
                charts.pressure.data.datasets[0].data = data.datasets.avg_pressure;
                charts.pressure.update();

                charts.humidity.data.labels = data.labels;
                charts.humidity.data.datasets[0].data = data.datasets.avg_humidity;
                charts.humidity.update();
            }
        } catch (error) {
            console.error('Error loading monthly daily stats:', error);
        }
    }

    async function loadGenericStats(type, updateCards = true) {
        // Just reuse the charts but fetch different data
        if (!charts.temperature) {
            // Re-init without multi if needed, but for now let's just use it
            charts.temperature = initChart('temperatureChart', i18n.temperature || 'Temperature', '#ef4444', '°C');
        }
        
        try {
            const response = await fetch(`/api/weather/${type}`);
            const data = await response.json();
            
            if (data.labels && data.datasets) {
                charts.temperature.data.labels = data.labels;
                charts.temperature.data.datasets[isMultiChart('temperatureChart') ? 1 : 0].data = data.datasets.avg_temperature;
                charts.temperature.update();
                if (charts.seaTemperature) {
                    charts.seaTemperature.data.labels = data.labels;
                    charts.seaTemperature.data.datasets[0].data = data.datasets.sea_temperature;
                    charts.seaTemperature.update();
                }

                if (charts.pressure) {
                    charts.pressure.data.labels = data.labels;
                    charts.pressure.data.datasets[0].data = data.datasets.avg_pressure;
                    charts.pressure.update();
                }

                if (charts.humidity) {
                    charts.humidity.data.labels = data.labels;
                    charts.humidity.data.datasets[0].data = data.datasets.avg_humidity;
                    charts.humidity.update();
                }
                
                if (updateCards) {
                    updateSummary(data.datasets);
                }
            }
        } catch (error) {
            console.error(`Error loading ${type} stats:`, error);
        }
    }

    async function loadMonthlySummary() {
        try {
            const response = await fetch('/api/weather/monthly-current');
            const data = await response.json();

            if (data.datasets) {
                updateSummary(data.datasets);
                updateMonthlySummaryTitle(data.through_date);
            }
        } catch (error) {
            console.error('Error loading current month summary:', error);
        }
    }
    
    function isMultiChart(id) {
         // Helper to check if chart was initialized as multi
         return id === 'temperatureChart' && isDailyLikePage;
    }

    async function loadSeaTemperatureHistoryForPath() {
        if (!charts.seaTemperature) return;

        let days = 30;
        if (path.includes('/statistics/recent')) days = 10;
        if (path.includes('/statistics/weekly')) days = 140;
        if (path.includes('/statistics/monthly')) days = 365;
        if (path.includes('/statistics/annual')) days = 730;

        try {
            const response = await fetch(`/api/water-temperatures/history?days=${encodeURIComponent(days)}&limit=2000`);
            const data = await response.json();
            if (!data || !Array.isArray(data.labels) || !Array.isArray(data.temperatures)) return;

            charts.seaTemperature.data.labels = data.labels;
            charts.seaTemperature.data.datasets[0].data = data.temperatures;
            charts.seaTemperature.update();
        } catch (error) {
            console.error('Error loading sea temperature history:', error);
        }
    }
});
