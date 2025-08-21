if (!sessionStorage.getItem("booted")) {
    let progress = 0;
  const progressBar = document.getElementById('booting-progress');
  const overlay = document.getElementById('booting-overlay');
  const logo = document.querySelector('img[alt="Majestic Coding Logo"]');
  const interval = setInterval(() => {
    progress = Math.min(progress + Math.random() * 8, 97);
    if (progressBar) progressBar.style.width = progress + "%";
  }, 120);
  
  window.addEventListener("DOMContentLoaded", function () {
    clearInterval(interval);
    if (progressBar) progressBar.style.width = "100%";
    setTimeout(() => {
      if (overlay) {
        overlay.style.opacity = "0";
        overlay.style.pointerEvents = "none";
        setTimeout(() => {
          overlay.style.display = "none";
          // Trigger bounce animation on logo
          if (logo) {
            logo.classList.add('logo-bounce-in');
          }
        }, 500);
      }
      sessionStorage.setItem("booted", "true");
    }, 600);
  });
} else {
  // Hide overlay immediately if not initial load
  const overlay = document.getElementById('booting-overlay');
  if (overlay) overlay.style.display = "none";
}