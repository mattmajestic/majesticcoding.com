async function updateLiveButton() {
  try {
    const res = await fetch("/api/stream/status");
    const isLive = await res.text();

    const btn = document.getElementById("live-btn");
    if (!btn) return;

    if (isLive.trim() === "true") {
      btn.classList.remove("bg-gray-700", "hover:bg-gray-600");
      btn.classList.add("bg-red-600", "hover:bg-red-500");
    } else {
      btn.classList.remove("bg-red-600", "hover:bg-red-500");
      btn.classList.add("bg-gray-700", "hover:bg-gray-600");
    }
  } catch (err) {
    console.error("Failed to fetch stream status:", err);
  }
}

// Run once on load
updateLiveButton();
// Poll every 15 seconds
setInterval(updateLiveButton, 15000);