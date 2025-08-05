 const chatMessages = document.getElementById('chat-messages');
const chatForm = document.getElementById('chat-form');
const chatInput = document.getElementById('chat-input');

const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws";
const wsHost = window.location.host;
const ws = new WebSocket(`${wsProtocol}://${wsHost}/ws/chat`);

function getColorForUsername(username) {
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 5) - hash);
  }
  const hue = hash % 360;
  return `hsl(${hue}, 70%, 60%)`; // colorful but readable
}


ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  const container = document.createElement('div');
  container.className = "mb-2";

  const meta = document.createElement('div');
  meta.className = "flex justify-between text-md text-gray-400";
  meta.innerHTML = `
    <span class="font-semibold" style="color: ${getColorForUsername(msg.Username)}">${msg.Username}</span>
    <span>${msg.Timestamp}</span>
    `;


  const content = document.createElement('div');
  content.className = "text-md text-gray-200";
  content.textContent = msg.Content;

  container.appendChild(meta);
  container.appendChild(content);
  chatMessages.appendChild(container);
  chatMessages.scrollTop = chatMessages.scrollHeight;
};

chatForm.addEventListener('submit', function(e) {
e.preventDefault();
const content = chatInput.value.trim();
if (!content) return;

const msg = { Content: content };
ws.send(JSON.stringify(msg));
chatInput.value = '';
});