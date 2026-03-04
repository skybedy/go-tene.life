document.addEventListener('DOMContentLoaded', function () {
    const i18n = window.pwsMapI18n || {};
    const mapEl = document.getElementById('pwsMap');
    const messageEl = document.getElementById('pwsMapMessage');
    const lastUpdateEl = document.getElementById('pwsLastUpdate');

    if (!mapEl) {
        return;
    }

    const map = L.map('pwsMap', {
        dragging: true,
        scrollWheelZoom: true,
        doubleClickZoom: true,
        boxZoom: false,
        keyboard: false,
        touchZoom: true,
        zoomControl: true
    }).setView([28.3, -15.8], 8);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 18,
        attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);
    const maxObservationAgeMs = 60 * 60 * 1000; // 1 hour
    function ensureLabelStyles() {
        if (document.getElementById('pwsTempLabelStyle')) return;
        const style = document.createElement('style');
        style.id = 'pwsTempLabelStyle';
        style.textContent = `
            .pws-temp-dot {
                width: 28px;
                height: 28px;
                border-radius: 9999px;
                border: 2px solid #1f2937;
                color: #fff;
                display: flex;
                align-items: center;
                justify-content: center;
                font-size: 11px;
                font-weight: 700;
                line-height: 1;
                box-shadow: 0 2px 8px rgba(0, 0, 0, 0.28);
                text-shadow: 0 1px 1px rgba(0, 0, 0, 0.4);
                user-select: none;
                white-space: nowrap;
            }
            .pws-temp-dot.is-na {
                color: #fff;
                font-size: 10px;
            }
            .pws-info-tooltip {
                background: rgba(17, 24, 39, 0.9);
                color: #fff;
                border: 1px solid rgba(255, 255, 255, 0.25);
                border-radius: 8px;
                padding: 6px 8px;
                font-size: 12px;
                line-height: 1.3;
                box-shadow: 0 2px 8px rgba(0, 0, 0, 0.25);
            }
            .pws-info-tooltip:before {
                border-top-color: rgba(17, 24, 39, 0.9);
            }
            @media (max-width: 640px) {
                .pws-temp-dot {
                    width: 24px;
                    height: 24px;
                    font-size: 10px;
                }
                .pws-temp-dot.is-na {
                    font-size: 9px;
                }
            }
        `;
        document.head.appendChild(style);
    }
    ensureLabelStyles();

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
            month: 'numeric',
            day: 'numeric',
            hour: 'numeric',
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

    function markerInfoText(point) {
        const name = point.name || point.stationId || '-';
        const lastUpdateLabel = i18n.lastUpdateLabel || 'Last update';
        return `<strong>${name}</strong><br>${lastUpdateLabel}: ${formatLocalDate(point.fetched_at_utc)}`;
    }

    function tempLabel(point) {
        if (point.temp_c == null || Number.isNaN(Number(point.temp_c))) {
            return '--';
        }
        return `${Math.round(Number(point.temp_c))}°`;
    }

    function markerIcon(point) {
        const label = tempLabel(point);
        const naClass = label === '--' ? ' is-na' : '';
        const bg = colorForTemp(point.temp_c);
        return L.divIcon({
            className: 'pws-temp-marker-wrap',
            html: `<div class="pws-temp-dot${naClass}" style="background:${bg}">${label}</div>`,
            iconSize: [28, 28],
            iconAnchor: [14, 14],
            popupAnchor: [0, -16]
        });
    }

    function parseISO(isoTime) {
        if (!isoTime) return null;
        const ts = new Date(isoTime);
        return Number.isNaN(ts.getTime()) ? null : ts;
    }

    function isRecentObservation(point) {
        const obs = parseISO(point.obs_time_utc);
        if (!obs) return false;
        return (Date.now() - obs.getTime()) <= maxObservationAgeMs;
    }

    function hasValidTemperature(point) {
        if (point.temp_c == null) return false;
        const temp = Number(point.temp_c);
        return Number.isFinite(temp);
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
                if (!isRecentObservation(point)) {
                    continue;
                }
                if (!hasValidTemperature(point)) {
                    continue;
                }

                if (point.fetched_at_utc) {
                    const currentFetched = parseISO(point.fetched_at_utc);
                    if (currentFetched && (!latestFetched || currentFetched > latestFetched)) {
                        latestFetched = currentFetched;
                    }
                }

                if (point.lat == null || point.lon == null) {
                    continue;
                }

                const marker = L.marker([point.lat, point.lon], {
                    icon: markerIcon(point)
                }).addTo(map);

                marker.bindPopup(popupText(point));
                marker.bindTooltip(markerInfoText(point), {
                    direction: 'top',
                    offset: [0, -24],
                    className: 'pws-info-tooltip'
                });
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
