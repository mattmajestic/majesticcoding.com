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
      .width(container.clientWidth || 800)
      .height(container.clientHeight || 600)
      .pointsData(points)
      .pointAltitude(d => 0.015)
      .pointColor(d => '#ff6b35')
      .pointRadius(d => 1.2)
      .pointLabel(d => {
        const timeAgo = d.time ? new Date(d.time).toLocaleDateString() : 'Recently';
        return `<b>${d.city}</b>, ${d.country}<br><small>Checked in: ${timeAgo}</small>`;
      })
      .labelsData(points)
      .labelLat(d => d.lat)
      .labelLng(d => d.lng)
      .labelText(d => d.city)
      .labelSize(1.8)
      .labelDotRadius(0.3)
      .labelColor(() => '#ffffff')
      .labelResolution(3)
      .labelAltitude(d => 0.02);

   
    globe.pointOfView({ lat: 45, lng: -30, altitude: 1.2 });


    globe.controls().enableZoom = true;
    globe.controls().enablePan = true;
    globe.controls().autoRotate = true;
    globe.controls().autoRotateSpeed = 0.8; // Slower for better viewing
    globe.controls().enableDamping = true;
    globe.controls().dampingFactor = 0.1;
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initGlobe);
  } else {
    initGlobe();
  }
})();