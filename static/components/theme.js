
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

let lastX, lastY;

document.addEventListener('mousemove', (e) => {
  const x = e.clientX;
  const y = e.clientY;

  if (lastX !== undefined && lastY !== undefined) {
    const dx = x - lastX;
    const dy = y - lastY;
    const distance = Math.sqrt(dx * dx + dy * dy);
    const steps = Math.max(1, Math.floor(distance / 4)); // 4px gap between dots

    for (let i = 1; i <= steps; i++) {
      const trailX = lastX + (dx * i) / steps;
      const trailY = lastY + (dy * i) / steps;
      createTrailDot(trailX, trailY);
    }
  } else {
    createTrailDot(x, y);
  }

  lastX = x;
  lastY = y;
});

function createTrailDot(x, y) {
  const dot = document.createElement('div');
  dot.className = 'cursor-trail';
  dot.style.left = `${x}px`;
  dot.style.top = `${y}px`;
  document.body.appendChild(dot);
  setTimeout(() => dot.remove(), 500); // Remove after fade
}