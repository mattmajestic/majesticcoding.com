// Docs Tab switching

document.querySelectorAll('.about-btn').forEach(button => {
    button.addEventListener('click', async () => {
      const tab = button.dataset.tab;
      const res = await fetch(`/about/${tab}`);
      const html = await res.text();
      document.getElementById('about-content').innerHTML = html;
    });
  });
