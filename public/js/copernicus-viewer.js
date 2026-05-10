(() => {
  const manifestUrl = '/data/copernicus/sea-temp/tenerife/2026/04/manifest.json';
  const statusEl = document.getElementById('copernicusStatus');
  const viewerEl = document.getElementById('copernicusViewer');
  if (!statusEl || !viewerEl) return;

  const imgEl = document.getElementById('copernicusFrame');
  const dateEl = document.getElementById('copernicusFrameDate');
  const indexEl = document.getElementById('copernicusFrameIndex');
  const sliderEl = document.getElementById('copernicusSlider');
  const playPauseEl = document.getElementById('copernicusPlayPause');
  const prevEl = document.getElementById('copernicusPrev');
  const nextEl = document.getElementById('copernicusNext');
  const speedEl = document.getElementById('copernicusSpeed');
  const videoWrapEl = document.getElementById('copernicusVideoWrap');
  const videoEl = document.getElementById('copernicusVideo');

  let manifest;
  let currentIndex = 0;
  let timer = null;

  const stopPlayback = () => {
    if (timer) clearInterval(timer);
    timer = null;
    playPauseEl.textContent = '▶ Play';
  };

  const preload = (index) => {
    const frame = manifest.frames[index];
    if (!frame) return;
    const img = new Image();
    img.src = new URL(frame.file, manifestUrl).toString();
  };

  const renderFrame = (index) => {
    const frame = manifest.frames[index];
    if (!frame) return;
    currentIndex = index;
    sliderEl.value = String(index);
    dateEl.textContent = frame.label || frame.date || '-';
    indexEl.textContent = `Snímek ${index + 1} / ${manifest.frames.length}`;
    imgEl.src = new URL(frame.file, manifestUrl).toString();
    preload(index + 1);
    preload(index - 1);
  };

  const step = (delta) => {
    const next = (currentIndex + delta + manifest.frames.length) % manifest.frames.length;
    renderFrame(next);
  };

  fetch(manifestUrl)
    .then((res) => {
      if (!res.ok) throw new Error(`Manifest nelze načíst (${res.status})`);
      return res.json();
    })
    .then((data) => {
      manifest = data;
      if (!manifest.frames || manifest.frames.length === 0) throw new Error('Manifest neobsahuje snímky.');

      sliderEl.max = String(manifest.frames.length - 1);
      statusEl.textContent = `${manifest.title || 'Satelitní vizualizace'} (${manifest.dateFrom} až ${manifest.dateTo})`;
      viewerEl.classList.remove('hidden');

      if (manifest.video) {
        videoWrapEl.classList.remove('hidden');
        videoEl.src = new URL(manifest.video, manifestUrl).toString();
      }

      renderFrame(0);

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
    })
    .catch((err) => {
      statusEl.textContent = `Nepodařilo se načíst Copernicus data: ${err.message}`;
    });
})();
