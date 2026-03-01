document.addEventListener('DOMContentLoaded', function () {
    const i18n = window.chartI18n || {};
    const charts = {
        temperature: null,
        pressure: null,
        humidity: null
    };

    function formatNumber(value, maxFractionDigits) {
        const numeric = Number(value);
        if (!Number.isFinite(numeric)) return value;
        return numeric.toLocaleString(undefined, {
            maximumFractionDigits: maxFractionDigits
        });
    }

    function initChart(id, label, color, unit, maxFractionDigits) {
        const ctx = document.getElementById(id).getContext('2d');
        return new Chart(ctx, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                    label: label,
                    data: [],
                    borderColor: color,
                    backgroundColor: color + '20',
                    borderWidth: 2,
                    fill: true,
                    tension: 0.4,
                    pointRadius: 2
                }]
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
                                return formatNumber(value, maxFractionDigits) + unit;
                            }
                        }
                    }
                },
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return `${label}: ${formatNumber(context.parsed.y, maxFractionDigits)}${unit}`;
                            }
                        }
                    }
                }
            }
        });
    }

    charts.temperature = initChart('temperatureChart', i18n.temperature || 'Temperature', '#ef4444', '°C', 1);
    charts.pressure = initChart('pressureChart', i18n.pressure || 'Pressure', '#3b82f6', ' hPa', 1);
    charts.humidity = initChart('humidityChart', i18n.humidity || 'Humidity', '#10b981', '%', 1);

    async function updateCharts() {
        try {
            const params = new URLSearchParams(window.location.search);
            const selectedDate = params.get('date');
            const endpoint = selectedDate
                ? `/api/weather/hourly?date=${encodeURIComponent(selectedDate)}`
                : '/api/weather/hourly';
            const response = await fetch(endpoint);
            const data = await response.json();

            if (data.labels && data.datasets) {
                // Temperature
                charts.temperature.data.labels = data.labels;
                charts.temperature.data.datasets[0].data = data.datasets.temperature;
                charts.temperature.update();

                // Pressure
                charts.pressure.data.labels = data.labels;
                charts.pressure.data.datasets[0].data = data.datasets.pressure;
                charts.pressure.update();

                // Humidity
                charts.humidity.data.labels = data.labels;
                charts.humidity.data.datasets[0].data = data.datasets.humidity;
                charts.humidity.update();
            }
        } catch (error) {
            console.error('Error updating charts:', error);
        }
    }

    updateCharts();
});
