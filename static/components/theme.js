
// Light/Dark Theme Toggle

document.addEventListener('DOMContentLoaded', () => {
  const body = document.body;
  const toggleBtn = document.getElementById('theme-toggle');
  const sunIcon = document.getElementById('icon-sun');
  const moonIcon = document.getElementById('icon-moon');

  const savedTheme = localStorage.getItem('theme');
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

  const applyTheme = (isDark) => {
    body.classList.toggle('dark', isDark);
    body.classList.toggle('light', !isDark);
    sunIcon.classList.toggle('hidden', isDark);
    moonIcon.classList.toggle('hidden', !isDark);
    localStorage.setItem('theme', isDark ? 'dark' : 'light');
  };

  if (savedTheme) {
    applyTheme(savedTheme === 'dark');
  } else {
    applyTheme(prefersDark);
  }

  if (toggleBtn) {
    toggleBtn.addEventListener('click', () => {
      const isDark = body.classList.contains('dark');
      applyTheme(!isDark);
    });
  }
});

document.addEventListener('DOMContentLoaded', () => {
  const trail = document.createElement('div');
  trail.className = 'cursor-trail glowing-effect';
  document.body.appendChild(trail);

  document.addEventListener('mousemove', (e) => {
    // Center the circle at the pointer
    trail.style.left = e.clientX + 'px';
    trail.style.top = e.clientY + 'px';
  });
});