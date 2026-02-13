document.addEventListener('DOMContentLoaded', function () {
    const charts = {
        temperature: null,
        pressure: null,
        humidity: null
    };

    function initChart(id, label, color, unit) {
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
                                return value + unit;
                            }
                        }
                    }
                },
                plugins: {
                    legend: { display: false }
                }
            }
        });
    }

    charts.temperature = initChart('temperatureChart', 'Teplota', '#ef4444', '°C');
    charts.pressure = initChart('pressureChart', 'Tlak', '#3b82f6', ' hPa');
    charts.humidity = initChart('humidityChart', 'Vlhkost', '#10b981', '%');

    async function updateCharts() {
        try {
            const response = await fetch('/api/weather/hourly');
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
