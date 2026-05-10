(() => {
  const statusEl = document.getElementById('copernicusStatus');
  const viewerEl = document.getElementById('copernicusViewer');
  const monthSelectEl = document.getElementById('copernicusMonthSelect');
  if (!viewerEl) return;
  const basePath = viewerEl.dataset.basePath || '/data/copernicus/sea-temp/tenerife';

  const imgEl = document.getElementById('copernicusFrame');
  const dateEl = document.getElementById('copernicusFrameDate');
  const indexEl = document.getElementById('copernicusFrameIndex');
  const sliderEl = document.getElementById('copernicusSlider');
  const playPauseEl = document.getElementById('copernicusPlayPause');
  const prevEl = document.getElementById('copernicusPrev');
  const nextEl = document.getElementById('copernicusNext');
  const speedEl = document.getElementById('copernicusSpeed');
  const labelPlay = viewerEl.dataset.labelPlay || 'Play';
  const labelPause = viewerEl.dataset.labelPause || 'Pause';
  const msgInvalidPeriod = viewerEl.dataset.msgInvalidPeriod || 'Invalid period.';
  const msgNoData = viewerEl.dataset.msgNoData || 'Data for this month is not available yet.';
  const msgManifestPrefix = viewerEl.dataset.msgManifestPrefix || 'Unable to load manifest';
  const msgLoadPrefix = viewerEl.dataset.msgLoadPrefix || 'Failed to load Copernicus data';
  const msgImageError = viewerEl.dataset.msgImageError || 'Error loading frame image.';

  let manifest;
  let manifestUrl = '';
  let manifestBaseUrl = '';
  let currentIndex = 0;
  let timer = null;

  const hideStatus = () => {
    if (!statusEl) return;
    statusEl.textContent = '';
    statusEl.classList.add('hidden');
  };

  const showStatus = (msg) => {
    if (!statusEl) return;
    statusEl.textContent = msg;
    statusEl.classList.remove('hidden');
  };

  const stopPlayback = () => {
    if (timer) clearInterval(timer);
    timer = null;
    playPauseEl.textContent = `▶ ${labelPlay}`;
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
      showStatus(msgInvalidPeriod);
      return;
    }

    const [year, mon] = parts;
    manifestUrl = `${basePath}/${year}/${mon}/manifest.json`;
    manifestBaseUrl = new URL(manifestUrl, window.location.origin).toString();
    hideStatus();
    hideViewer();

    fetch(manifestBaseUrl)
      .then((res) => {
        if (!res.ok) {
          if (res.status === 404) {
            throw new Error(msgNoData);
          }
          throw new Error(`${msgManifestPrefix} (${res.status})`);
        }
        return res.json();
      })
      .then((data) => {
        manifest = data;
        if (!manifest.frames || manifest.frames.length === 0) {
          throw new Error(msgNoData);
        }

        sliderEl.max = String(manifest.frames.length - 1);
        hideStatus();
        viewerEl.classList.remove('hidden');

        renderFrame(0);
      })
      .catch((err) => {
        hideViewer();
        showStatus(err.message.includes(msgNoData)
          ? err.message
          : `${msgLoadPrefix}: ${err.message}`);
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
    playPauseEl.textContent = `⏸ ${labelPause}`;
    timer = setInterval(() => step(1), Number(speedEl.value));
  });
  speedEl.addEventListener('change', () => {
    if (!timer) return;
    stopPlayback();
    playPauseEl.click();
  });

  imgEl.addEventListener('error', () => {
    showStatus(msgImageError);
    stopPlayback();
  });

  if (monthSelectEl) {
    monthSelectEl.addEventListener('change', () => loadMonth(monthSelectEl.value));
    loadMonth(monthSelectEl.value);
    return;
  }

  loadMonth('2026-04');
})();
