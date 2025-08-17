(function () {
  const log = document.getElementById('log');

  // Build WS URL (supports ?ws=... and optional ?room=...)
  const q = new URLSearchParams(location.search);
  const base = (location.protocol === 'https:' ? 'wss' : 'ws') + '://' + location.host;
  const roomPart = q.get('room') ? ('?room=' + encodeURIComponent(q.get('room'))) : '';
  const wsURL = q.get('ws') || (base + '/ws/chat' + roomPart);

  function getColorForUsername(u) {
    let h = 0;
    for (let i = 0; i < (u || '').length; i++) h = u.charCodeAt(i) + ((h << 5) - h);
    const hue = ((h % 360) + 360) % 360;
    return `hsl(${hue},70%,60%)`;
  }

  function append({ Username, Content } = {}) {
    const row = document.createElement('div');
    row.className = 'cw-msg';

    if (Username) {
      const name = document.createElement('div');     // block element = its own line
      name.className = 'cw-name';
      name.textContent = Username;
      name.style.color = getColorForUsername(Username); // dynamic per-user color
      row.appendChild(name);
    }

    const text = document.createElement('div');       // block element = next line
    text.className = 'cw-text';
    text.textContent = Content ?? '';
    row.appendChild(text);

    log.appendChild(row);
    log.scrollTop = log.scrollHeight; // stick to bottom
  }

  const ws = new WebSocket(wsURL);
  ws.onmessage = (e) => {
    try { append(JSON.parse(e.data)); }
    catch { append({ Content: String(e.data) }); }
  };
})();
