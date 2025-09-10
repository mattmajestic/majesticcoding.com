(function () {
  const overlay = document.getElementById('overlay');
  if (!overlay) return;

  // container styles
  Object.assign(overlay.style, {
    position: 'fixed',
    top: '10%',
    left: '50%',
    transform: 'translateX(-50%)',
    display: 'flex',
    flexDirection: 'row',
    gap: '32px',
    pointerEvents: 'none',
    fontFamily: 'ui-sans-serif, system-ui, sans-serif',
    zIndex: '999999',
  });

  // allow-list (match either username OR display_name, case-insensitive)
  const allowed = new Set([
    'majesticcodingtwitch',
    'pungentgurgi',
    'bigdaddddddy69',
    'bitcoin__',
  ]);

  // resolve image by either username or display_name
  function imgForName(username, displayName) {
    const k1 = (username || '').toLowerCase();
    const k2 = (displayName || '').toLowerCase();
    const table = {
      'majesticcodingtwitch': '/static/img/arsenal.gif',
      'pungentgurgi': '/static/img/manu.gif',
      'bigdaddddddy69': '/static/img/spurs.gif',
      'bitcoin__': '/static/img/bitcoin.gif',
    };
    return table[k1] || table[k2] || '/static/img/arsenal-2.gif';
  }

  function makeChip(displayName, message, username) {
    // chip wrapper
    const wrap = document.createElement('div');
    Object.assign(wrap.style, {
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      opacity: '0',
      transition: 'opacity 180ms ease',
    });

    // bubble
    const bubble = document.createElement('div');
    bubble.textContent = message || '';
    Object.assign(bubble.style, {
      background: 'rgba(0,0,0,0.65)',
      color: '#fff',
      padding: '8px 12px',
      borderRadius: '12px',
      border: '2px solid #31b3ff',
      whiteSpace: 'nowrap',
      maxWidth: '220px',
      overflow: 'hidden',
      textOverflow: 'ellipsis',
      fontSize: '14px',
      marginBottom: '8px',
    });
    wrap.appendChild(bubble);

    // image
    const img = document.createElement('img');
    img.src = imgForName(username, displayName);
    Object.assign(img.style, {
      height: '60px',
      filter: 'drop-shadow(0 2px 6px rgba(0,0,0,.5))',
      marginBottom: '6px',
    });
    wrap.appendChild(img);

    // display name pill
    const name = document.createElement('div');
    name.textContent = displayName || username || 'Unknown';
    Object.assign(name.style, {
      color: '#fff',
      fontWeight: '600',
      background: 'rgba(0,0,0,0.55)',
      border: '2px solid #fff',
      borderRadius: '9999px',
      padding: '2px 8px',
      fontSize: '12px',
    });
    wrap.appendChild(name);

    // insert + animate
    overlay.appendChild(wrap);
    requestAnimationFrame(() => { wrap.style.opacity = '1'; });

    // auto-remove after 8s
    setTimeout(() => {
      wrap.style.opacity = '0';
      setTimeout(() => wrap.remove(), 220);
    }, 8000);
  }

  // WebSocket: /ws/twitch
  const ws = new WebSocket(
    (location.protocol === 'https:' ? 'wss' : 'ws') + '://' + location.host + '/ws/twitch'
  );

  ws.onmessage = (ev) => {
    try {
      const raw = JSON.parse(ev.data);
      const items = Array.isArray(raw) ? raw : [raw];

      items.forEach((msg) => {
        const dn = (msg.display_name || '').toLowerCase();
        const un = (msg.username || '').toLowerCase();

        // only show if in allow-list
        if (!allowed.has(un) && !allowed.has(dn)) return;

        const displayName = msg.display_name || msg.username || '';
        const username = msg.username || msg.display_name || '';
        const message = msg.message || '';

        makeChip(displayName, message, username);
      });
    } catch (e) {
      console.warn('bad message', e);
    }
  };

  ws.onerror = (e) => console.error('ws error', e);

  // --- optional debug smoke test ---
  // setTimeout(() => makeChip('MajesticCodingTwitch', '⚽ Arsenal fan here!', 'MajesticCodingTwitch'), 1500);
})();
