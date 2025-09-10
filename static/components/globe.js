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

  // Initialize globe with JSON data
  async function initGlobe() {
    try {
      await waitForGlobe();
      console.log('Globe library loaded');
    } catch (e) {
      console.error('Failed to load Globe library:', e);
      return;
    }

    // Load recent checkins from API
    let points = [];
    try {
      const response = await fetch('/api/checkins/recent');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const checkins = await response.json();
      console.log('Loaded recent checkins:', checkins);
      
      points = checkins.map(checkin => ({
        lat: checkin.lat,
        lng: checkin.lon, // Note: DB uses 'lon' instead of 'lng'
        size: 0.8,
        color: 'orange',
        city: checkin.city || `${checkin.lat.toFixed(2)}, ${checkin.lon.toFixed(2)}`, // Use city name or coordinates
        country: checkin.country || 'Recent Checkin',
        time: checkin.checkin_time
      }));
      
      console.log('Mapped points:', points);
    } catch (error) {
      console.error('Failed to load recent checkins:', error);
      return;
    }

    // Initialize globe
    const container = document.getElementById('globe');
    if (!container) {
      console.error('Globe container not found');
      return;
    }

    globe = Globe()(container)
      .globeImageUrl('//unpkg.com/three-globe/example/img/earth-dark.jpg')
      .backgroundColor('#000913')
      .width(container.clientWidth || 800)
      .height(container.clientHeight || 600)
      .pointsData(points)
      .pointAltitude(d => 0.02)
      .pointColor(d => 'red')
      .pointRadius(d => 1.0)
      .pointLabel(d => {
        const timeAgo = d.time ? new Date(d.time).toLocaleDateString() : 'Recently';
        return `<b>${d.city}</b>, ${d.country}<br><small>Checked in: ${timeAgo}</small>`;
      })
      .labelsData(points)
      .labelLat(d => d.lat)
      .labelLng(d => d.lng)
      .labelText(d => d.city)
      .labelSize(1.5)
      .labelDotRadius(0.5)
      .labelColor(() => 'white')
      .labelResolution(2)
      .labelAltitude(d => 0.025);

    // Enable controls
    globe.controls().enableZoom = true;
    globe.controls().autoRotate = true;
    globe.controls().autoRotateSpeed = 0.5;
    
    console.log('Globe initialized with', points.length, 'points');
    
    // Debug: Log first point details
    if (points.length > 0) {
      console.log('First point:', points[0]);
    }
    
    // Test: Add a few seconds delay then log the globe's point data
    setTimeout(() => {
      console.log('Globe points data after init:', globe.pointsData());
    }, 2000);
  }

  // Start when page loads
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initGlobe);
  } else {
    initGlobe();
  }
})();