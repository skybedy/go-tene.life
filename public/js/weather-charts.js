document.addEventListener('DOMContentLoaded', function () {
    const i18n = window.chartI18n || {};
    const hourlyTableBody = document.getElementById('hourlyHomeTableBody');
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
        const canvas = document.getElementById(id);
        if (!canvas || !window.Chart) return null;
        const ctx = canvas.getContext('2d');
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

    function renderHourlyTable(labels, datasets) {
        if (!hourlyTableBody) return;

        if (!Array.isArray(labels) || !labels.length) {
            hourlyTableBody.innerHTML = '<tr><td colspan="4" class="px-3 py-8 text-center text-gray-500 italic">' + (i18n.noData || 'No data available.') + '</td></tr>';
            return;
        }

        function formatValue(value, digits, unit) {
            return typeof value === 'number' && Number.isFinite(value)
                ? value.toFixed(digits) + unit
                : '--';
        }

        hourlyTableBody.innerHTML = labels.map(function(_, index) {
            const reverseIndex = labels.length - 1 - index;
            const label = labels[reverseIndex];
            const temperature = datasets && Array.isArray(datasets.temperature) ? datasets.temperature[reverseIndex] : null;
            const pressure = datasets && Array.isArray(datasets.pressure) ? datasets.pressure[reverseIndex] : null;
            const humidity = datasets && Array.isArray(datasets.humidity) ? datasets.humidity[reverseIndex] : null;

            return '<tr class="odd:bg-white/55 even:bg-blue-50/55 hover:bg-white/75 transition">'
                + '<td class="px-3 py-3 font-medium text-gray-900">' + label + '</td>'
                + '<td class="px-3 py-3 text-right text-gray-700">' + formatValue(temperature, 1, ' °C') + '</td>'
                + '<td class="px-3 py-3 text-right text-gray-700">' + formatValue(pressure, 1, ' hPa') + '</td>'
                + '<td class="px-3 py-3 text-right text-gray-700">' + formatValue(humidity, 0, ' %') + '</td>'
                + '</tr>';
        }).join('');
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
                renderHourlyTable(data.labels, data.datasets);

                if (charts.temperature) {
                    charts.temperature.data.labels = data.labels;
                    charts.temperature.data.datasets[0].data = data.datasets.temperature;
                    charts.temperature.update();
                }

                if (charts.pressure) {
                    charts.pressure.data.labels = data.labels;
                    charts.pressure.data.datasets[0].data = data.datasets.pressure;
                    charts.pressure.update();
                }

                if (charts.humidity) {
                    charts.humidity.data.labels = data.labels;
                    charts.humidity.data.datasets[0].data = data.datasets.humidity;
                    charts.humidity.update();
                }
            }
        } catch (error) {
            console.error('Error updating charts:', error);
            renderHourlyTable([], null);
        }
    }

    updateCharts();
});
