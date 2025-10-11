// Simple Globe with JSON checkins
(function() {
  let globe;

  // Wait for Globe library to load
  function waitForGlobe() {
    return new Promise((resolve, reject) => {
      if (typeof Globe !== 'undefined') {
        resolve();
        return;
      }
      
      let attempts = 0;
      const maxAttempts = 50;
      const checkInterval = setInterval(() => {
        attempts++;
        if (typeof Globe !== 'undefined') {
          clearInterval(checkInterval);
          resolve();
        } else if (attempts >= maxAttempts) {
          clearInterval(checkInterval);
          reject(new Error('Globe library failed to load'));
        }
      }, 100);
    });
  }

  async function initGlobe() {
    try {
      await waitForGlobe();
      console.log('Globe library loaded');
    } catch (e) {
      console.error('Failed to load Globe library:', e);
      return;
    }

    let points = [];
    try {
      const response = await fetch('/api/checkins/recent');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const checkins = await response.json();
      
      if (checkins && Array.isArray(checkins)) {
        points = checkins.map(checkin => ({
          lat: checkin.lat,
          lng: checkin.lon,
          size: 0.8,
          color: 'orange',
          username: checkin.username || 'mattmajestic',
          city: checkin.city || 'Recent Checkin',
          country: checkin.country || 'Unknown',
          time: checkin.checkin_time
        }));
      } else {
        console.warn('No checkins data available or invalid format');
        points = []; // Use empty array as fallback
      }
    } catch (error) {
      console.error('Failed to load recent checkins:', error);
      points = [];
    }

    const container = document.getElementById('globe');
    if (!container) {
      console.error('Globe container not found');
      return;
    }

    globe = Globe()(container)
      .globeImageUrl('//unpkg.com/three-globe/example/img/earth-blue-marble.jpg')
      .backgroundColor('rgba(0,0,0,0)')
      .atmosphereColor('#00ffff')
      .atmosphereAltitude(0.25)
      .width(container.clientWidth || 800)
      .height(container.clientHeight || 600)
      .pointsData(points)
      .pointAltitude(d => 0.01)
      .pointColor(d => '#00ffff')
      .pointRadius(d => 2.5)
      .pointsMerge(true)
      .pointLabel(d => {
        const timeAgo = d.time ? new Date(d.time).toLocaleDateString() : 'Recently';
        return `
          <div style="background: linear-gradient(135deg, rgba(0,10,20,0.95), rgba(0,30,60,0.95)); backdrop-filter: blur(15px); padding: 14px 18px; border-radius: 12px; color: white; font-family: 'Courier New', 'Consolas', monospace; text-align: center; min-width: 180px; border: 2px solid rgba(0,255,255,0.6); box-shadow: 0 0 20px rgba(0,255,255,0.4), inset 0 0 20px rgba(0,255,255,0.1); position: relative;">
            <div style="position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, transparent, #00ffff, transparent); animation: scan 2s linear infinite;"></div>
            <div style="font-weight: 700; font-size: 16px; color: #00ffff; margin-bottom: 8px; text-shadow: 0 0 10px rgba(0,255,255,0.8); letter-spacing: 1px;">${d.city.toUpperCase()}</div>
            <div style="font-size: 13px; color: #66ffff; margin-bottom: 6px; opacity: 0.9;">${d.country}</div>
            <div style="font-size: 11px; color: #00ff88; margin-top: 8px; padding-top: 8px; border-top: 1px solid rgba(0,255,255,0.3); font-family: 'Courier New', monospace;">⚡ SIGNAL: ${timeAgo}</div>
          </div>
          <style>
            @keyframes scan {
              0%, 100% { transform: translateY(0); opacity: 0; }
              50% { transform: translateY(20px); opacity: 1; }
            }
          </style>
        `;
      })
      .labelsData(points)
      .labelLat(d => d.lat)
      .labelLng(d => d.lng)
      .labelText(d => d.city)
      .labelSize(4)
      .labelDotRadius(1.2)
      .labelColor(() => '#00ffff')
      .labelResolution(6)
      .labelAltitude(d => 0.12);

    // Animate points with pulsing effect
    setInterval(() => {
      globe.pointsData(points.map(p => ({
        ...p,
        size: 0.8 + Math.random() * 0.4
      })));
    }, 1000);

    globe.pointOfView({ lat: 30, lng: -20, altitude: 2.5 });

    globe.controls().enableZoom = true;
    globe.controls().enablePan = true;
    globe.controls().autoRotate = true;
    globe.controls().autoRotateSpeed = 0.8;
    globe.controls().enableDamping = true;
    globe.controls().dampingFactor = 0.1;
    globe.controls().minDistance = 180;
    globe.controls().maxDistance = 800;
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initGlobe);
  } else {
    initGlobe();
  }
})();