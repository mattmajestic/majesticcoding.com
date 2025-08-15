document.addEventListener('DOMContentLoaded', async () => {
      // Fetch stream status
      const streamEl = document.getElementById('stream-status');
      try {
        const res = await fetch('/api/stream/status');
        const data = await res.json();
        streamEl.textContent = data.status ? 'Live' : 'Offline';
      } catch (e) {
        streamEl.textContent = 'Error';
      }

      // Fetch User Count
      const chatUsers = document.getElementById('chat-users');
      try {
        const res = await fetch('/api/chat/users');
        const data = await res.json();
        chatUsers.textContent = data.user_count !== 'undefined' ? data.user_count : 'N/A';
      } catch (e) {
        chatUsers.textContent = 'Error';
      }

    // Fetch and display socials as JSON
    const providers = ['github', 'youtube', 'twitch', 'leetcode'];
    const socials = {};
    for (const provider of providers) {
      try {
        const res = await fetch(`/api/stats/${provider}`);
        socials[provider] = await res.json();
      } catch (e) {
        socials[provider] = "Error loading";
      }
    }
    document.getElementById('provider-json').textContent = JSON.stringify(socials, null, 2);
    document.getElementById('socials-spinner-container').style.display = 'none';

      // System metrics
      try {
        const res = await fetch('/api/metrics');
        const metrics = await res.json();
        const metTbody = document.getElementById('system-metrics');
        Object.entries(metrics).forEach(([key, value]) => {
          let displayValue = value;
          if (typeof value === 'number' && !Number.isInteger(value)) {
            displayValue = value.toFixed(2);
          }
          const row = document.createElement('tr');
          row.innerHTML = `
            <td class="border px-4 py-2">${key}</td>
            <td class="border px-4 py-2">${displayValue}</td>
          `;
          metTbody.appendChild(row);
        });
        document.getElementById('metrics-spinner-container').style.display = 'none';
      } catch (e) {
        const row = document.createElement('tr');
        row.innerHTML = `<td class="border px-4 py-2 text-red-500" colspan="2">Error loading metrics</td>`;
        document.getElementById('system-metrics').appendChild(row);
        document.getElementById('metrics-spinner-container').style.display = 'none';
      }
        try {
            const res = await fetch('/api/git/hash'); // or your new route, e.g. /api/git/latest
            const data = await res.json();
            document.getElementById('git-hash').textContent = `Date: ${data.commit_date}`;
            document.getElementById('git-message').textContent = `Message: ${data.message}`;
        } catch (e) {
            document.getElementById('git-hash').textContent = 'Error loading commit date';
            document.getElementById('git-message').textContent = 'Error loading commit message';
        }
    });

document.addEventListener('DOMContentLoaded', () => {
  const btn = document.getElementById('metrics-collapse-btn');
  const content = document.getElementById('metrics-content');
  const icon = document.getElementById('metrics-collapse-icon');
  if (btn && content && icon) {
    btn.addEventListener('click', () => {
      content.classList.toggle('hidden');
      icon.classList.toggle('fa-chevron-up');
      icon.classList.toggle('fa-chevron-down');
    });
  }
});