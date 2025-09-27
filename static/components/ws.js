// Websockets appending Chat Messages
// AI Commands included as dummy for now

// Chat Components from chat.tmpl
const chatMessages = document.getElementById('chat-messages');
const chatForm = document.getElementById('chat-form');
const chatInput = document.getElementById('chat-input');

// Websocket Setup
const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws";
const wsHost = window.location.host;

let ws;

// Function to create WebSocket connection with auth
function createWebSocketConnection() {
  if (ws) {
    ws.close();
  }

  // Get auth token for WebSocket connection
  let wsUrl = `${wsProtocol}://${wsHost}/ws/chat`;
  const token = localStorage.getItem('supabase_token');
  if (token) {
    console.log('🔐 Connecting to chat with auth token');
    wsUrl += `?token=${encodeURIComponent(token)}`;
  } else {
    console.log('👤 Connecting to chat as anonymous user');
  }

  ws = new WebSocket(wsUrl);
  setupWebSocketHandlers();
}

// Color of Username
function getColorForUsername(username) {
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 5) - hash);
  }
  const hue = hash % 360;

  // Check if light mode
  const isLightMode = document.body.classList.contains('light');

  if (isLightMode) {
    return `hsl(${hue}, 70%, 30%)`; // darker colors for light mode
  } else {
    return `hsl(${hue}, 70%, 60%)`; // original colors for dark mode
  }
}

// Setup WebSocket event handlers
function setupWebSocketHandlers() {
  // Websocket Connection
  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);

    const container = document.createElement('div');
    container.className = "mb-2";

    const meta = document.createElement('div');
    meta.className = "flex justify-between text-md chat-meta-text";
    meta.innerHTML = `
      <span class="font-semibold" style="color: ${getColorForUsername(msg.Username)}">${msg.Username}</span>
      <span>${msg.DisplayTime}</span>
      `;

    const content = document.createElement('div');
    content.className = "text-md chat-content-text";
    content.textContent = msg.Content;

    container.appendChild(meta);
    container.appendChild(content);
    chatMessages.appendChild(container);
    chatMessages.scrollTop = chatMessages.scrollHeight;
  };
}

// Form submission handler (outside of WebSocket setup)
if (chatForm) {
  chatForm.addEventListener('submit', function (e) {
    e.preventDefault();
    const content = chatInput.value.trim();
    if (!content) return;

    // Check if user is authenticated
    const token = localStorage.getItem('supabase_token');
    if (!token) {
      // Don't submit, let the auth UI handle showing login prompt
      return;
    }

    const msg = { Content: content };
    ws.send(JSON.stringify(msg));

  if (content.startsWith('!ai ')) {
    const prompt = content.slice(4);

    // Get auth token - try localStorage first (more reliable)
    let token = localStorage.getItem('supabase_token');

    // Fallback to Supabase session if available
    if (!token && window.supabaseAuth) {
      token = window.supabaseAuth.getAuthToken();
    }

    if (!token) {
      const aiMsg = {
        Content: "🤖 Please log in to use AI chat",
        Username: "AI"
      };
      ws.send(JSON.stringify(aiMsg));
      return;
    }

    // Call your AI API
    fetch('/api/llm/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        prompt: prompt,
        provider: 'gemini'
      })
    })
    .then(response => response.json())
    .then(data => {
      const aiMsg = {
        Content: "🤖 " + data.response,
        Username: "AI (" + data.provider + ")"
      };
      ws.send(JSON.stringify(aiMsg));
    })
    .catch(error => {
      console.error('AI API error:', error);
      const aiMsg = {
        Content: "🤖 Sorry, AI is temporarily unavailable",
        Username: "AI"
      };
      ws.send(JSON.stringify(aiMsg));
    });
  }


    chatInput.value = '';
  });
}

function appendMessage(msg) {
  const container = document.createElement('div');
  container.className = "mb-2";

  const meta = document.createElement('div');
  meta.className = "flex justify-between text-md chat-meta-text";
  meta.innerHTML = `
    <span class="font-semibold" style="color: ${getColorForUsername(msg.Username)}">${msg.Username}</span>
    <span>${msg.DisplayTime}</span>
  `;

  const content = document.createElement('div');
  content.className = "text-md chat-content-text";
  content.textContent = msg.Content;

  container.appendChild(meta);
  container.appendChild(content);
  chatMessages.appendChild(container);
  chatMessages.scrollTop = chatMessages.scrollHeight;
}

// Function to update chat UI based on auth state
function updateChatAuthUI() {
  const token = localStorage.getItem('supabase_token');
  const authMessage = document.getElementById('auth-status-message');
  const chatForm = document.getElementById('chat-form');

  console.log('🔐 Updating chat auth UI, token:', token ? 'present' : 'missing');

  if (!token) {
    // User not authenticated - show auth message, hide form
    if (authMessage) {
      authMessage.classList.remove('hidden');
      console.log('📝 Showing auth status message');
    }
    if (chatForm) {
      chatForm.style.display = 'none';
      console.log('🚫 Hiding chat form');
    }
  } else {
    // User authenticated - hide auth message, show form
    if (authMessage) {
      authMessage.classList.add('hidden');
      console.log('✅ Hiding auth status message');
    }
    if (chatForm) {
      chatForm.style.display = 'flex';
      console.log('💬 Showing chat form');
    }
  }
}

// Function to reconnect WebSocket when auth state changes
function reconnectChat() {
  console.log('🔄 Reconnecting chat with new auth state...');
  chatInitialized = false; // Reset flag to allow reinitialization
  initializeChat();
  updateChatAuthUI(); // Update UI immediately
}

// Initialize chat after a short delay to allow auth to initialize
let authInitialized = false;
let chatInitialized = false;

function initializeChat() {
  if (chatInitialized) return;
  chatInitialized = true;
  createWebSocketConnection();
  updateChatAuthUI();

  // Keep updating auth UI every 500ms until we have a token or auth manager
  const authCheckInterval = setInterval(() => {
    const token = localStorage.getItem('supabase_token');
    const authManager = window.authManager;

    updateChatAuthUI();

    // Stop checking once we have auth state established
    if (token || (authManager && authManager.supabase)) {
      clearInterval(authCheckInterval);
      console.log('🔐 Auth state established, stopping auth UI updates');
    }
  }, 500);
}

// Check if auth is ready, otherwise wait for it
function checkAuthAndInitialize() {
  // Check if we have auth token or if auth manager is ready
  const token = localStorage.getItem('supabase_token');
  const authManager = window.authManager;

  if (token || (authManager && authManager.getCurrentUser())) {
    // We have auth info, initialize chat
    initializeChat();
  } else if (authManager && authManager.supabase) {
    // Auth manager is initialized but no user - initialize chat anyway for anonymous users
    initializeChat();
  } else {
    // Wait a bit more for auth to initialize
    setTimeout(checkAuthAndInitialize, 100);
  }
}

// Start the initialization check
setTimeout(checkAuthAndInitialize, 50);

// Listen for auth state changes to reconnect chat
window.addEventListener('storage', (e) => {
  if (e.key === 'supabase_token') {
    reconnectChat();
  }
});

// Also expose reconnection function globally for auth manager
window.reconnectChat = reconnectChat;
