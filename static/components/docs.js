// Docs Tab switching with URL support

function loadTab(tabName) {
  // For tab switching, we need to navigate to the new URL
  // This will trigger a full page load with the correct content
  window.location.href = `/docs/${tabName}`;
}

// Initialize active tab button based on current URL
function initializeTab() {
  const currentPath = window.location.pathname;
  const validTabs = ['api', 'partials', 'chat', 'database', 'stream-docs', 'hosting', 'ai', 'cache'];

  // Extract section from URL like /docs/partials
  const pathParts = currentPath.split('/');
  const section = pathParts[2]; // /docs/section-name

  if (section && validTabs.includes(section)) {
    // We're on a specific docs page - just update the active button
    // Content is already loaded by the server
    updateActiveButton(section);
  } else if (currentPath === '/docs' || currentPath === '/docs/') {
    // Default to first tab if we're on the base docs page
    updateActiveButton('api');
  }
}

function updateActiveButton(activeTab) {
  document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.classList.remove('bg-blue-600');
    btn.classList.add('bg-gray-800');
  });

  const activeBtn = document.querySelector(`[data-tab="${activeTab}"]`);
  if (activeBtn) {
    activeBtn.classList.remove('bg-gray-800');
    activeBtn.classList.add('bg-blue-600');
  }
}

// Handle tab button clicks
document.querySelectorAll('.tab-btn').forEach(button => {
    button.addEventListener('click', () => {
      const tab = button.dataset.tab;

      if (tab && !button.onclick) {
        loadTab(tab);
      }
      // If it has onclick, let the inline handler take care of it
    });
});

// Initialize on page load
document.addEventListener('DOMContentLoaded', initializeTab);

function showTab(tab) {
  document.getElementById('tab-gcp').classList.add('hidden');
  document.getElementById('tab-aws').classList.add('hidden');
  document.getElementById('tab-' + tab).classList.remove('hidden');
}

// Global tab functions for sub-sections
window.showStreamTab = function(tabName) {
  document.querySelectorAll('.stream-tab-content').forEach(tab => {
    tab.classList.add('hidden');
  });
  document.getElementById('stream-tab-' + tabName).classList.remove('hidden');

  // Update button states - remove active styling from all stream buttons
  document.querySelectorAll('.stream-tab-btn').forEach(btn => {
    btn.classList.remove('bg-gray-700', 'text-white');
    btn.classList.add('text-blue-300');
  });

  // Find and highlight the active button
  const activeButton = document.querySelector(`button[onclick="showStreamTab('${tabName}')"]`);
  if (activeButton) {
    activeButton.classList.remove('text-blue-300');
    activeButton.classList.add('bg-gray-700', 'text-white');
  }
}

window.showDbTab = function(tabName) {
  document.querySelectorAll('.db-tab-content').forEach(tab => {
    tab.classList.add('hidden');
  });
  document.getElementById('db-tab-' + tabName).classList.remove('hidden');

  // Update button states - remove active styling from all db buttons
  document.querySelectorAll('.db-tab-btn').forEach(btn => {
    btn.classList.remove('bg-gray-700', 'text-white');
    btn.classList.add('text-blue-300');
  });

  // Find and highlight the active button
  const activeButton = document.querySelector(`button[onclick="showDbTab('${tabName}')"]`);
  if (activeButton) {
    activeButton.classList.remove('text-blue-300');
    activeButton.classList.add('bg-gray-700', 'text-white');
  }
}

window.showChatTab = function(tabName) {
  // Hide all chat tab content
  document.querySelectorAll('.chat-tab-content').forEach(tab => {
    tab.classList.add('hidden');
  });
  document.getElementById('chat-tab-' + tabName).classList.remove('hidden');

  // Update button states - remove active styling from all chat buttons
  document.querySelectorAll('.chat-tab-btn').forEach(btn => {
    btn.classList.remove('bg-gray-700', 'text-white');
    btn.classList.add('text-blue-300');
  });

  // Find and highlight the active button
  const activeButton = document.querySelector(`button[onclick="showChatTab('${tabName}')"]`);
  if (activeButton) {
    activeButton.classList.remove('text-blue-300');
    activeButton.classList.add('bg-gray-700', 'text-white');
  }
}

window.showCacheTab = function(tabName) {
  document.querySelectorAll('.cache-tab-content').forEach(tab => {
    tab.classList.add('hidden');
  });
  document.getElementById('cache-tab-' + tabName).classList.remove('hidden');

  // Update button states - remove active styling from all cache buttons
  document.querySelectorAll('.cache-tab-btn').forEach(btn => {
    btn.classList.remove('bg-gray-700', 'text-white');
    btn.classList.add('text-blue-300');
  });

  // Find and highlight the active button
  const activeButton = document.querySelector(`button[onclick="showCacheTab('${tabName}')"]`);
  if (activeButton) {
    activeButton.classList.remove('text-blue-300');
    activeButton.classList.add('bg-gray-700', 'text-white');
  }
}

// Add new function for partials tab switching
window.showPartialsTab = function(tabName) {
  document.querySelectorAll('.partials-tab-content').forEach(tab => {
    tab.classList.add('hidden');
  });
  document.getElementById('partials-tab-' + tabName).classList.remove('hidden');

  // Update button states - remove active styling from all partials buttons
  document.querySelectorAll('.partials-tab-btn').forEach(btn => {
    btn.classList.remove('bg-gray-700', 'text-white');
    btn.classList.add('text-blue-300');
  });

  // Find and highlight the active button
  const activeButton = document.querySelector(`button[onclick="showPartialsTab('${tabName}')"]`);
  if (activeButton) {
    activeButton.classList.remove('text-blue-300');
    activeButton.classList.add('bg-gray-700', 'text-white');
  }
}

// Add new function for hosting tab switching
window.showHostingTab = function(tabName) {
  document.querySelectorAll('.hosting-tab-content').forEach(tab => {
    tab.classList.add('hidden');
  });
  document.getElementById('hosting-tab-' + tabName).classList.remove('hidden');

  // Update button states - remove active styling from all hosting buttons
  document.querySelectorAll('.hosting-tab-btn').forEach(btn => {
    btn.classList.remove('bg-gray-700', 'text-white');
    btn.classList.add('text-blue-300');
  });

  // Find and highlight the active button
  const activeButton = document.querySelector(`button[onclick="showHostingTab('${tabName}')"]`);
  if (activeButton) {
    activeButton.classList.remove('text-blue-300');
    activeButton.classList.add('bg-gray-700', 'text-white');
  }
}

// Add new function for API tab switching
window.showApiTab = function(tabName) {
  document.querySelectorAll('.api-tab-content').forEach(tab => {
    tab.classList.add('hidden');
  });
  document.getElementById('api-tab-' + tabName).classList.remove('hidden');

  // Update button states - remove active styling from all API buttons
  document.querySelectorAll('.api-tab-btn').forEach(btn => {
    btn.classList.remove('bg-gray-700', 'text-white');
    btn.classList.add('text-blue-300');
  });

  // Find and highlight the active button
  const activeButton = document.querySelector(`button[onclick="showApiTab('${tabName}')"]`);
  if (activeButton) {
    activeButton.classList.remove('text-blue-300');
    activeButton.classList.add('bg-gray-700', 'text-white');
  }
}

// Add new function for AI tab switching
window.showAiTab = function(tabName) {
  document.querySelectorAll('.ai-tab-content').forEach(tab => {
    tab.classList.add('hidden');
  });
  document.getElementById('ai-tab-' + tabName).classList.remove('hidden');

  // Update button states - remove active styling from all AI buttons
  document.querySelectorAll('.ai-tab-btn').forEach(btn => {
    btn.classList.remove('bg-gray-700', 'text-white');
    btn.classList.add('text-blue-300');
  });

  // Find and highlight the active button
  const activeButton = document.querySelector(`button[onclick="showAiTab('${tabName}')"]`);
  if (activeButton) {
    activeButton.classList.remove('text-blue-300');
    activeButton.classList.add('bg-gray-700', 'text-white');
  }
}