document.addEventListener('DOMContentLoaded', function () {
    // This is a simplified version of statistics charts
    // In a real app, you would have different charts for each page

    async function loadStats(type) {
        try {
            const response = await fetch(`/api/weather/${type}`);
            const data = await response.json();
            console.log(`Loaded ${type} stats:`, data);
            // Here you would initialize Chart.js for tables/charts
            // For now, we are focusing on the tables that are already in the HTML
        } catch (error) {
            console.error(`Error loading ${type} stats:`, error);
        }
    }

    const path = window.location.pathname;
    if (path.includes('/statistics/daily')) loadStats('daily');
    if (path.includes('/statistics/weekly')) loadStats('weekly');
    if (path.includes('/statistics/monthly')) loadStats('monthly');
    if (path.includes('/statistics/annual')) loadStats('annual');
});
