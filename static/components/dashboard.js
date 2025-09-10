document.addEventListener('DOMContentLoaded', async () => {
  // Update stream status and indicator
  const streamEl = document.getElementById('stream-status');
  const streamIndicator = document.getElementById('stream-indicator');
  try {
    const res = await fetch('/api/stream/status');
    const data = await res.json();
    const isLive = data.status;
    streamEl.textContent = isLive ? 'Live' : 'Offline';
    
    if (streamIndicator) {
      const dot = streamIndicator.querySelector('.status-dot');
      if (dot) {
        dot.classList.remove('bg-red-500', 'bg-green-500');
        dot.classList.add(isLive ? 'bg-green-500' : 'bg-red-500');
      }
    }
  } catch (e) {
    streamEl.textContent = 'Error';
  }

  // Update chat users count
  const chatUsers = document.getElementById('chat-users');
  try {
    const res = await fetch('/api/chat/users');
    const data = await res.json();
    chatUsers.textContent = data.user_count !== 'undefined' ? data.user_count : '0';
  } catch (e) {
    chatUsers.textContent = 'N/A';
  }

  // Update git information
  try {
    const res = await fetch('/api/git/hash');
    const data = await res.json();
    document.getElementById('git-hash').innerHTML = `<span class="text-xs text-gray-400">${data.commit_date}</span>`;
    document.getElementById('git-message').innerHTML = `${data.message}`;
  } catch (e) {
    document.getElementById('git-hash').innerHTML = '<span class="text-red-400">Error loading</span>';
    document.getElementById('git-message').innerHTML = 'Error loading commit';
  }

  // Update system health (CPU, Memory, and Uptime)
  updateSystemHealth();

  // Load metrics cards
  await loadMetricsCards();
});

// Update system health bars
async function updateSystemHealth() {
  try {
    const res = await fetch('/api/metrics');
    const metrics = await res.json();
    
    // Update CPU usage
    const cpuUsage = metrics.cpu_usage || Math.random() * 100; // Fallback to random for demo
    const cpuBar = document.getElementById('cpu-usage');
    const cpuValue = document.getElementById('cpu-value');
    if (cpuBar && cpuValue) {
      const roundedCpu = Math.round(cpuUsage * 100) / 100;
      cpuBar.style.width = `${roundedCpu}%`;
      cpuValue.textContent = `${roundedCpu}%`;
      
      // Update color based on usage
      cpuBar.classList.remove('bg-green-500', 'bg-yellow-500', 'bg-red-500');
      if (roundedCpu < 50) cpuBar.classList.add('bg-green-500');
      else if (roundedCpu < 80) cpuBar.classList.add('bg-yellow-500');
      else cpuBar.classList.add('bg-red-500');
    }
    
    // Update Memory usage
    const memoryUsage = metrics.memory_usage || Math.random() * 100; // Fallback to random for demo
    const memoryBar = document.getElementById('memory-usage');
    const memoryValue = document.getElementById('memory-value');
    if (memoryBar && memoryValue) {
      const roundedMemory = Math.round(memoryUsage * 100) / 100;
      memoryBar.style.width = `${roundedMemory}%`;
      memoryValue.textContent = `${roundedMemory}%`;
      
      // Update color based on usage
      memoryBar.classList.remove('bg-green-500', 'bg-yellow-500', 'bg-red-500');
      if (roundedMemory < 50) memoryBar.classList.add('bg-green-500');
      else if (roundedMemory < 80) memoryBar.classList.add('bg-yellow-500');
      else memoryBar.classList.add('bg-red-500');
    }
    
    // Update Uptime
    const uptimeEl = document.getElementById('uptime-value');
    if (uptimeEl) {
      const uptime = metrics.uptime || Date.now() / 1000; // Fallback to current time in seconds
      uptimeEl.textContent = formatUptime(uptime);
    }
  } catch (e) {
    console.error('Error updating system health:', e);
  }
}

// Load and populate metrics cards
async function loadMetricsCards() {
  const container = document.getElementById('metrics-container');
  if (!container) return;
  
  try {
    const res = await fetch('/api/metrics');
    const metrics = await res.json();
    
    // Clear loading state
    container.innerHTML = '';
    
    // Create metric cards
    Object.entries(metrics).forEach(([key, value]) => {
      const card = createMetricCard(key, value);
      container.appendChild(card);
    });
    
  } catch (e) {
    container.innerHTML = `
      <div class="metrics-error">
        <i class="fas fa-exclamation-triangle text-red-400"></i>
        <span>Error loading metrics</span>
      </div>
    `;
  }
}

// Create individual metric card
function createMetricCard(key, value) {
  const card = document.createElement('div');
  card.className = 'metric-card';
  
  let displayValue = value;
  
  // Format different metric types
  if (typeof value === 'number') {
    if (key.includes('memory') || key.includes('disk')) {
      displayValue = formatBytes(value);
    } else if (key.includes('cpu')) {
      // Handle very small CPU values with scientific notation
      if (value < 0.01) {
        displayValue = value.toExponential(2);
      } else {
        displayValue = `${Math.round(value * 100) / 100}%`;
      }
    } else if (key.includes('usage')) {
      displayValue = `${Math.round(value * 100) / 100}%`;
    } else if (!Number.isInteger(value)) {
      displayValue = Math.round(value * 100) / 100;
    }
  }
  
  // Get appropriate icon for metric
  const icon = getMetricIcon(key);
  
  card.innerHTML = `
    <div class="metric-header">
      <i class="${icon} metric-icon"></i>
      <span class="metric-name">${formatMetricName(key)}</span>
    </div>
    <div class="metric-value">${displayValue}</div>
  `;
  
  return card;
}

// Format metric names for display
function formatMetricName(key) {
  return key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
}

// Get icon for metric type
function getMetricIcon(key) {
  if (key.includes('cpu')) return 'fas fa-microchip';
  if (key.includes('memory')) return 'fas fa-memory';
  if (key.includes('disk')) return 'fas fa-hdd';
  if (key.includes('network')) return 'fas fa-network-wired';
  if (key.includes('time') || key.includes('duration')) return 'fas fa-clock';
  if (key.includes('count') || key.includes('total')) return 'fas fa-hashtag';
  return 'fas fa-chart-line';
}

// Format bytes for display
function formatBytes(bytes) {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Format uptime for display
function formatUptime(seconds) {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  
  if (days > 0) return `${days}d ${hours}h`;
  if (hours > 0) return `${hours}h ${minutes}m`;
  return `${minutes}m`;
}

// Refresh metrics function
function refreshMetrics() {
  const container = document.getElementById('metrics-container');
  if (container) {
    container.innerHTML = `
      <div class="metrics-loading">
        <i class="fas fa-spinner fa-spin"></i>
        <span>Refreshing metrics...</span>
      </div>
    `;
  }
  
  // Reload all data
  updateSystemHealth();
  loadMetricsCards();
}