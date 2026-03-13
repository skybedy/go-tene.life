document.addEventListener('DOMContentLoaded', function () {
    const i18n = window.statsI18n || {};
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
                                return value + unit;
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
    charts.temperature = initChart('temperatureChart', i18n.temperature || 'Temperature', '#ef4444', '°C', true); // Multi for min/avg/max
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

                if (charts.seaTemperature && data.datasets.sea_temperature) {
                    charts.seaTemperature.data.labels = data.labels;
                    charts.seaTemperature.data.datasets[0].data = data.datasets.sea_temperature;
                    charts.seaTemperature.update();
                }

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

        const avg = arr => arr.length ? (arr.reduce((a, b) => a + b, 0) / arr.length).toFixed(1) : '--';

        const tempAvg = avg(datasets.avg_temperature);
        const pressAvg = avg(datasets.avg_pressure);
        const humAvg = avg(datasets.avg_humidity);

        if (document.getElementById('stat-temp-avg')) document.getElementById('stat-temp-avg').textContent = tempAvg + ' °C';
        if (document.getElementById('stat-pressure-avg')) document.getElementById('stat-pressure-avg').textContent = pressAvg + ' hPa';
        if (document.getElementById('stat-humidity-avg')) document.getElementById('stat-humidity-avg').textContent = Math.round(humAvg) + ' %';
    }

    // Route handling
    const path = window.location.pathname;
    if (path.includes('/statistics/daily')) {
        loadDailyStats();
    } else if (path.includes('/statistics/recent')) {
        loadDailyStats();
    } else if (path.includes('/statistics/weekly')) {
        // Implement weekly...
        loadGenericStats('weekly');
    } else if (path.includes('/statistics/monthly')) {
        loadGenericStats('monthly');
    } else if (path.includes('/statistics/annual')) {
        loadGenericStats('annual');
    }

    async function loadGenericStats(type) {
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

                if (charts.seaTemperature && data.datasets.sea_temperature) {
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
                
                updateSummary(data.datasets);
            }
        } catch (error) {
            console.error(`Error loading ${type} stats:`, error);
        }
    }
    
    function isMultiChart(id) {
         // Helper to check if chart was initialized as multi
         return id === 'temperatureChart' && path.includes('/statistics/daily');
    }
});
