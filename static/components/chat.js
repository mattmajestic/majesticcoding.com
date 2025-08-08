// Chat Sends

const chatToggle = document.getElementById('chat-toggle');
const chatWidget = document.getElementById('chat-widget');
const hiddenMsg = document.getElementById('chat-hidden-msg');
const eyeOpen = document.getElementById('eye-open');
const eyeClosed = document.getElementById('eye-closed');

chatToggle.addEventListener('click', () => {
  chatWidget.classList.toggle('collapsed');
  hiddenMsg.classList.toggle('hidden');
  eyeOpen.classList.toggle('hidden');
  eyeClosed.classList.toggle('hidden');
});

// Update the User Counts
function updateUserCount() {
  fetch('/api/chat/users')
    .then(res => res.json())
    .then(data => {
      const count = data.user_count;
      document.getElementById('user-count').textContent = count;
      document.getElementById('user-count-label').textContent = count === 1 ? 'Chatter' : 'Chatters';
    });
}

document.getElementById('chat-input').addEventListener('keydown', function(e) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault(); // Prevent newline
    document.getElementById('chat-form').dispatchEvent(new Event('submit', {cancelable: true, bubbles: true}));
  }
});

// Insert AI Command into Chat Input
document.getElementById('insert-ai').addEventListener('click', () => {
  const input = document.getElementById('chat-input');
  input.value = '!ai ';
  input.focus();
});


updateUserCount();
setInterval(updateUserCount, 10000);
