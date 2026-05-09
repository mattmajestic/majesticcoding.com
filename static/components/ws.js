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
  const token = getValidClerkToken();

  try {
    const protocols = token ? ['clerk-auth', token] : undefined;
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

// Returns clerk_token only if it's actually a Clerk JWT (iss contains 'clerk').
// Clears and returns null if it's a stale token from another provider (e.g. Supabase).
function getValidClerkToken() {
  const token = localStorage.getItem('clerk_token');
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')));
    if (!payload.iss || !payload.iss.includes('clerk')) {
      localStorage.removeItem('clerk_token');
      return null;
    }
  } catch (_) { /* malformed JWT — let server reject it */ }
  return token;
}

// Simple function to check if user is authenticated
function isUserAuthenticated() {
  return getValidClerkToken() !== null;
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
  // Force close existing connection so we reconnect with the new auth state
  if (ws) {
    ws.close();
    ws = null;
  }
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
  // Wait for clerk-auth.js to finish loading and caching the token
  if (window.clerkAuthReady) {
    initializeChat();
  } else {
    document.addEventListener('clerk-auth-ready', initializeChat, { once: true });
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
  if (e.key === 'clerk_token') {
    reconnectChat();
  }
});
document.addEventListener('clerk-auth-change', reconnectChat);

// Also expose reconnection function globally for auth manager
window.reconnectChat = reconnectChat;
