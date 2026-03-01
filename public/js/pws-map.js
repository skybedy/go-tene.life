document.addEventListener('DOMContentLoaded', function () {
    const i18n = window.pwsMapI18n || {};
    const mapEl = document.getElementById('pwsMap');
    const messageEl = document.getElementById('pwsMapMessage');
    const lastUpdateEl = document.getElementById('pwsLastUpdate');

    if (!mapEl) {
        return;
    }

    const map = L.map('pwsMap').setView([28.2916, -16.6291], 10);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 18,
        attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    function colorForTemp(temp) {
        if (temp == null || Number.isNaN(temp)) return '#6b7280';
        if (temp < 12) return '#2563eb';
        if (temp < 18) return '#0ea5e9';
        if (temp < 24) return '#10b981';
        if (temp < 30) return '#f59e0b';
        return '#ef4444';
    }

    function formatLocalDate(isoTime) {
        if (!isoTime) return '-';
        const ts = new Date(isoTime);
        if (Number.isNaN(ts.getTime())) return '-';

        return new Intl.DateTimeFormat(document.documentElement.lang || undefined, {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        }).format(ts);
    }

    function setLastUpdate(latestISO) {
        const label = i18n.lastUpdateLabel || 'Last update';
        lastUpdateEl.textContent = `${label}: ${formatLocalDate(latestISO)}`;
    }

    function showMessage(text) {
        if (!messageEl) return;
        messageEl.textContent = text;
        messageEl.classList.remove('hidden');
    }

    function popupText(point) {
        const tempText = point.temp_c == null ? '--' : `${Number(point.temp_c).toFixed(1)} °C`;
        const lines = [
            `<strong>${point.name || point.stationId}</strong>`,
            `${tempText}`,
            `${i18n.observedAt || 'Observed'}: ${formatLocalDate(point.obs_time_utc)}`,
            `${i18n.fetchedAt || 'Fetched'}: ${formatLocalDate(point.fetched_at_utc)}`
        ];

        if (point.stale) {
            lines.push(`<em>${i18n.stale || 'Stale'}</em>`);
        }
        if (point.invalid) {
            lines.push(`<em>${i18n.invalid || 'Invalid'}</em>`);
        }
        return lines.join('<br>');
    }

    async function loadPoints() {
        try {
            const response = await fetch('/api/tenerife/pws-latest');
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }

            const points = await response.json();
            if (!Array.isArray(points) || points.length === 0) {
                setLastUpdate(null);
                showMessage(i18n.noData || 'No station data available.');
                return;
            }

            let latestFetched = null;
            let visiblePoints = 0;

            for (const point of points) {
                if (point.fetched_at_utc) {
                    const currentFetched = new Date(point.fetched_at_utc);
                    if (!Number.isNaN(currentFetched.getTime()) && (!latestFetched || currentFetched > latestFetched)) {
                        latestFetched = currentFetched;
                    }
                }

                if (point.lat == null || point.lon == null) {
                    continue;
                }

                const marker = L.circleMarker([point.lat, point.lon], {
                    radius: 8,
                    weight: 2,
                    color: '#1f2937',
                    fillColor: colorForTemp(point.temp_c),
                    fillOpacity: 0.85
                }).addTo(map);

                marker.bindPopup(popupText(point));
                visiblePoints++;
            }

            if (visiblePoints === 0) {
                showMessage(i18n.noData || 'No station data available.');
            }

            setLastUpdate(latestFetched ? latestFetched.toISOString() : null);
        } catch (error) {
            console.error('Failed to load PWS map points:', error);
            setLastUpdate(null);
            showMessage(i18n.noData || 'No station data available.');
        }
    }

    loadPoints();
});
