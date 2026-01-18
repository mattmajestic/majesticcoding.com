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
let wsConnecting = false;

// Function to create WebSocket connection with auth
function createWebSocketConnection() {
  // Prevent multiple simultaneous connection attempts
  if (wsConnecting) {
    return;
  }

  // Close existing connection if open
  if (ws && ws.readyState !== WebSocket.CLOSED) {
    ws.close();
  }

  wsConnecting = true;

  // Get auth token for WebSocket connection
  const wsUrl = `${wsProtocol}://${wsHost}/ws/chat`;
  const token = localStorage.getItem('supabase_token');

  try {
    const protocols = token ? ['supabase-auth', token] : undefined;
    ws = new WebSocket(wsUrl, protocols);
    setupWebSocketHandlers();
  } catch (error) {
    wsConnecting = false;
  }
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
  // Connection opened
  ws.onopen = () => {
    wsConnecting = false;
  };

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

  // Handle WebSocket errors silently
  ws.onerror = (error) => {
    wsConnecting = false;
  };

  // Handle WebSocket close silently
  ws.onclose = (event) => {
    wsConnecting = false;
  };
}

// Form submission handler (outside of WebSocket setup)
if (chatForm) {
  chatForm.addEventListener('submit', function (e) {
    e.preventDefault();
    const content = chatInput.value.trim();
    if (!content) return;

    // Check if user is authenticated
    if (!isUserAuthenticated()) {
      // Redirect to auth page
      window.location.href = '/auth';
      return;
    }

    const msg = { Content: content };
    ws.send(JSON.stringify(msg));

  if (content.startsWith('!ai ')) {
    // AI feature disabled - show coming soon message
    const aiMsg = {
      Content: "🤖 AI chat coming soon! Currently under maintenance.",
      Username: "AI"
    };
    ws.send(JSON.stringify(aiMsg));
    chatInput.value = ''; // Clear input after AI command
    return;
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

// Simple function to check if user is authenticated
function isUserAuthenticated() {
  return localStorage.getItem('supabase_token') !== null;
}

// Function to update chat UI based on auth state
function updateChatAuthUI() {
  const authMessage = document.getElementById('auth-status-message');
  const chatForm = document.getElementById('chat-form');

  if (!isUserAuthenticated()) {
    // Show "sign in required" message, hide form
    if (authMessage) authMessage.classList.remove('hidden');
    if (chatForm) chatForm.style.display = 'none';
  } else {
    // Hide message, show form
    if (authMessage) authMessage.classList.add('hidden');
    if (chatForm) chatForm.style.display = 'flex';
  }
}

// Function to reconnect WebSocket when auth state changes
function reconnectChat() {
  // Only reconnect if WebSocket is not already open
  if (ws && ws.readyState === WebSocket.OPEN) {
    // Already connected, just update UI
    updateChatAuthUI();
    return;
  }

  // Reset and reconnect
  chatInitialized = false;
  wsConnecting = false;
  initializeChat();
  updateChatAuthUI();
}

// Initialize chat
let chatInitialized = false;

function initializeChat() {
  if (chatInitialized) return;
  chatInitialized = true;
  createWebSocketConnection();
  updateChatAuthUI();
}

// Wait for page to fully load and auth to be ready before initializing chat
function waitForAuthAndInitialize() {
  // Check if auth manager is available and initialized
  if (window.authManager && window.authManager.supabase) {
    // Auth is ready, wait a bit more to ensure any OAuth flow is complete
    setTimeout(initializeChat, 1000);
  } else {
    // Auth not ready yet, check again in 200ms
    setTimeout(waitForAuthAndInitialize, 200);
  }
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', waitForAuthAndInitialize);
} else {
  // Document already loaded
  waitForAuthAndInitialize();
}

// Listen for auth state changes to reconnect chat
window.addEventListener('storage', (e) => {
  if (e.key === 'supabase_token') {
    reconnectChat();
  }
});

// Also expose reconnection function globally for auth manager
window.reconnectChat = reconnectChat;
