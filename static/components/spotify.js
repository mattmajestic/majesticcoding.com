// Spotify Terminal Widget - Updates existing HTML elements
(function() {
  let spotifyData = null;
  let currentProgress = 0;
  let isPlaying = false;
  let progressInterval = null;
  let lastTrackId = null;
  
  // DOM elements
  const elements = {
    status: document.getElementById('status'),
    title: document.getElementById('track-title'),
    artist: document.getElementById('track-artist'),
    progressBar: document.getElementById('progress-bar'),
    timeDisplay: document.getElementById('time-display'),
    albumImage: document.getElementById('album-image')
  };
  
  function formatTime(ms) {
    const minutes = Math.floor(ms / 60000);
    const seconds = Math.floor((ms % 60000) / 1000);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  }
  
  function createProgressBar(progress, total) {
    const percentage = progress / total;
    const filled = Math.round(percentage * 20);
    
    let progressHtml = '';
    for (let i = 0; i < 20; i++) {
      if (i < filled) {
        progressHtml += '<span class="progress-filled">#</span>';
      } else {
        progressHtml += '<span class="progress-empty">#</span>';
      }
    }
    
    return progressHtml;
  }
  
  function updateDisplay() {
    if (!spotifyData) return;
    
    // Update track info
    elements.title.textContent = `"${spotifyData.title}"`;
    elements.artist.textContent = spotifyData.artist;
    
    // Update album image
    if (spotifyData.albumImage) {
      elements.albumImage.src = spotifyData.albumImage;
      elements.albumImage.alt = `${spotifyData.album} album artwork`;
    }
    
    // Update progress bar and time
    const progressHtml = createProgressBar(currentProgress, spotifyData.durationMs);
    const currentTime = formatTime(currentProgress);
    const totalTime = formatTime(spotifyData.durationMs);
    
    elements.progressBar.innerHTML = progressHtml;
    elements.timeDisplay.textContent = `${currentTime}/${totalTime}`;
    elements.status.textContent = isPlaying ? 'playing' : 'paused';
  }
  
  function updateProgress() {
    if (!spotifyData || !isPlaying) return;
    
    // Increment progress by 1 second
    currentProgress += 1000;
    
    // Loop back to start if we reach the end
    if (currentProgress >= spotifyData.durationMs) {
      currentProgress = 0;
    }
    
    // Update just the progress bar and time
    const progressHtml = createProgressBar(currentProgress, spotifyData.durationMs);
    const currentTime = formatTime(currentProgress);
    const totalTime = formatTime(spotifyData.durationMs);
    
    elements.progressBar.innerHTML = progressHtml;
    elements.timeDisplay.textContent = `${currentTime}/${totalTime}`;
  }
  
  function startProgress() {
    if (isPlaying && !progressInterval) {
      progressInterval = setInterval(updateProgress, 1000);
    }
  }
  
  function stopProgress() {
    if (progressInterval) {
      clearInterval(progressInterval);
      progressInterval = null;
    }
  }
  
  async function loadSpotifyData(showLoading = true) {
    try {
      if (showLoading) {
        elements.status.textContent = 'loading...';
      }
      
      // Try real Spotify API first
      try {
        const response = await fetch('/api/spotify/current');
        
        if (response.ok) {
          // Real Spotify data from your account
          const data = await response.json();
          
          // Check if track changed (only update if different)
          const trackId = data.url || data.title + data.artist;
          const trackChanged = trackId !== lastTrackId;
          
          let newSpotifyData;
          
          // Handle case where no song is playing
          if (!data.is_playing && data.message) {
            newSpotifyData = {
              title: "No Track Playing",
              artist: "No Artist Playing",
              album: "Not Active",
              albumImage: "https://avatars.githubusercontent.com/u/33904170?v=4",
              isPlaying: false,
              progressMs: 0,
              durationMs: 0,
              trackId: null
            };
          } else {
            // Map API response to expected format
            newSpotifyData = {
              title: data.title,
              artist: data.artists ? data.artists.join(', ') : 'Unknown',
              album: data.album,
              albumImage: data.album_image,
              isPlaying: data.is_playing,
              progressMs: data.progress_ms,
              durationMs: data.duration_ms,
              trackId: trackId
            };
          }
          
          // Only update if track changed or play state changed
          if (trackChanged || !spotifyData || spotifyData.isPlaying !== newSpotifyData.isPlaying) {
            spotifyData = newSpotifyData;
            lastTrackId = trackId;
            console.log('Track changed, updating display:', spotifyData);
            
            currentProgress = spotifyData.progressMs;
            isPlaying = spotifyData.isPlaying;
            
            updateDisplay();
            
            // Manage progress interval
            if (isPlaying) {
              startProgress();
            } else {
              stopProgress();
            }
          } else {
            // Just update progress for current track
            currentProgress = newSpotifyData.progressMs;
          }
        } else {
          // API not available, use sample data
          throw new Error('API not available');
        }
      } catch (apiError) {
        // Fall back to sample JSON file (only for initial load)
        if (showLoading) {
          console.log('Using sample data:', apiError.message);
          const response = await fetch('/static/data/spotify.json');
          spotifyData = await response.json();
          console.log('Loaded sample data:', spotifyData);
          
          currentProgress = spotifyData.progressMs;
          isPlaying = spotifyData.isPlaying;
          
          updateDisplay();
          
          if (isPlaying) {
            startProgress();
          }
        }
      }
      
    } catch (error) {
      console.error('Error loading Spotify data:', error);
      elements.status.textContent = 'error';
      elements.title.textContent = '"Error loading data"';
      elements.artist.textContent = 'Unknown';
    }
  }
  
  // Auto-refresh every minute to sync with Spotify
  function startAutoRefresh() {
    setInterval(() => loadSpotifyData(false), 60000); // Poll every minute without showing loading
  }
  
  // Initialize when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      loadSpotifyData();
      startAutoRefresh();
    });
  } else {
    loadSpotifyData();
    startAutoRefresh();
  }
})();