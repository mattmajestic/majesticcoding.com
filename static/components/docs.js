// Docs Tab switching

document.querySelectorAll('.tab-btn').forEach(button => {
    button.addEventListener('click', async () => {
      const tab = button.dataset.tab;
      const res = await fetch(`/docs/${tab}`);
      const html = await res.text();
      document.getElementById('tab-content').innerHTML = html;
    });
  });

function showTab(tab) {
  document.getElementById('tab-gcp').classList.add('hidden');
  document.getElementById('tab-aws').classList.add('hidden');
  document.getElementById('tab-' + tab).classList.remove('hidden');
}