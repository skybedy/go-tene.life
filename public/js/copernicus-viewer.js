(() => {
  const basePath = '/data/copernicus/sea-temp/tenerife';
  const statusEl = document.getElementById('copernicusStatus');
  const viewerEl = document.getElementById('copernicusViewer');
  const monthSelectEl = document.getElementById('copernicusMonthSelect');
  if (!statusEl || !viewerEl) return;

  const imgEl = document.getElementById('copernicusFrame');
  const dateEl = document.getElementById('copernicusFrameDate');
  const indexEl = document.getElementById('copernicusFrameIndex');
  const sliderEl = document.getElementById('copernicusSlider');
  const playPauseEl = document.getElementById('copernicusPlayPause');
  const prevEl = document.getElementById('copernicusPrev');
  const nextEl = document.getElementById('copernicusNext');
  const speedEl = document.getElementById('copernicusSpeed');

  let manifest;
  let manifestUrl = '';
  let manifestBaseUrl = '';
  let currentIndex = 0;
  let timer = null;

  const stopPlayback = () => {
    if (timer) clearInterval(timer);
    timer = null;
    playPauseEl.textContent = '▶ Play';
  };

  const hideViewer = () => {
    stopPlayback();
    viewerEl.classList.add('hidden');
    imgEl.removeAttribute('src');
    dateEl.textContent = '-';
    indexEl.textContent = '-';
    sliderEl.value = '0';
    sliderEl.max = '0';
  };

  const preload = (index) => {
    if (!manifest || index < 0 || index >= manifest.frames.length) return;
    const frame = manifest.frames[index];
    if (!frame) return;
    const img = new Image();
    img.src = new URL(frame.file, manifestBaseUrl).toString();
  };

  const renderFrame = (index) => {
    const frame = manifest.frames[index];
    if (!frame) return;
    currentIndex = index;
    sliderEl.value = String(index);
    dateEl.textContent = frame.label || frame.date || '-';
    indexEl.textContent = `Snímek ${index + 1} / ${manifest.frames.length}`;
    imgEl.src = new URL(frame.file, manifestBaseUrl).toString();
    preload(index + 1);
    preload(index - 1);
  };

  const step = (delta) => {
    const next = (currentIndex + delta + manifest.frames.length) % manifest.frames.length;
    renderFrame(next);
  };

  const loadMonth = (month) => {
    const parts = month.split('-');
    if (parts.length !== 2) {
      hideViewer();
      statusEl.textContent = 'Neplatné období.';
      return;
    }

    const [year, mon] = parts;
    manifestUrl = `${basePath}/${year}/${mon}/manifest.json`;
    manifestBaseUrl = new URL(manifestUrl, window.location.origin).toString();
    statusEl.textContent = `Načítám manifest pro ${month}…`;
    hideViewer();

    fetch(manifestBaseUrl)
      .then((res) => {
        if (!res.ok) {
          if (res.status === 404) {
            throw new Error('Data pro tento měsíc zatím nejsou dostupná.');
          }
          throw new Error(`Manifest nelze načíst (${res.status})`);
        }
        return res.json();
      })
      .then((data) => {
        manifest = data;
        if (!manifest.frames || manifest.frames.length === 0) {
          throw new Error('Data pro tento měsíc zatím nejsou dostupná.');
        }

        sliderEl.max = String(manifest.frames.length - 1);
        statusEl.textContent = `${manifest.title || 'Satelitní vizualizace'} (${manifest.dateFrom} až ${manifest.dateTo})`;
        viewerEl.classList.remove('hidden');

        renderFrame(0);
      })
      .catch((err) => {
        hideViewer();
        statusEl.textContent = err.message.includes('Data pro tento měsíc')
          ? err.message
          : `Nepodařilo se načíst Copernicus data: ${err.message}`;
      });
  };

  sliderEl.addEventListener('input', () => {
    stopPlayback();
    renderFrame(Number(sliderEl.value));
  });
  prevEl.addEventListener('click', () => { stopPlayback(); step(-1); });
  nextEl.addEventListener('click', () => { stopPlayback(); step(1); });
  playPauseEl.addEventListener('click', () => {
    if (timer) {
      stopPlayback();
      return;
    }
    playPauseEl.textContent = '⏸ Pause';
    timer = setInterval(() => step(1), Number(speedEl.value));
  });
  speedEl.addEventListener('change', () => {
    if (!timer) return;
    stopPlayback();
    playPauseEl.click();
  });

  imgEl.addEventListener('error', () => {
    statusEl.textContent = 'Chyba při načítání snímku.';
    stopPlayback();
  });

  if (monthSelectEl) {
    monthSelectEl.addEventListener('change', () => loadMonth(monthSelectEl.value));
    loadMonth(monthSelectEl.value);
    return;
  }

  loadMonth('2026-04');
})();
