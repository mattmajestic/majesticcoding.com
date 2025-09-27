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
      .globeImageUrl('//unpkg.com/three-globe/example/img/earth-night.jpg')
      .backgroundColor('rgba(0,0,0,0)')
      .atmosphereColor('rgba(255,107,53,0.6)')
      .atmosphereAltitude(0.15)
      .width(container.clientWidth || 800)
      .height(container.clientHeight || 600)
      .pointsData(points)
      .pointAltitude(d => 0.02)
      .pointColor(d => '#ff6b35')
      .pointRadius(d => 1.8)
      .pointsMerge(true)
      .pointLabel(d => {
        const timeAgo = d.time ? new Date(d.time).toLocaleDateString() : 'Recently';
        return `
          <div style="background: linear-gradient(135deg, rgba(0,0,0,0.9), rgba(30,30,50,0.9)); backdrop-filter: blur(10px); padding: 12px; border-radius: 8px; color: white; font-family: 'Inter', 'Segoe UI', sans-serif; text-align: center; min-width: 140px; border: 1px solid rgba(255,107,53,0.3); box-shadow: 0 8px 32px rgba(0,0,0,0.3);">
            <div style="font-weight: 600; font-size: 15px; color: #ff6b35; margin-bottom: 6px; text-shadow: 0 2px 4px rgba(0,0,0,0.5);">${d.city}</div>
            <div style="font-size: 12px; color: #e0e0e0; margin-bottom: 4px;">${d.country}</div>
            <div style="font-size: 10px; color: #a0a0a0; margin-top: 6px; padding-top: 6px; border-top: 1px solid rgba(255,107,53,0.2);">📍 ${timeAgo}</div>
          </div>
        `;
      })
      .labelsData(points)
      .labelLat(d => d.lat)
      .labelLng(d => d.lng)
      .labelText(d => d.city)
      .labelSize(3.5)
      .labelDotRadius(0.8)
      .labelColor(() => '#ff6b35')
      .labelResolution(6)
      .labelAltitude(d => 0.08);

   
    globe.pointOfView({ lat: 30, lng: -20, altitude: 2.2 });


    globe.controls().enableZoom = true;
    globe.controls().enablePan = true;
    globe.controls().autoRotate = true;
    globe.controls().autoRotateSpeed = 0.5;
    globe.controls().enableDamping = true;
    globe.controls().dampingFactor = 0.15;
    globe.controls().minDistance = 180;
    globe.controls().maxDistance = 800;
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initGlobe);
  } else {
    initGlobe();
  }
})();