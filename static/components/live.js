async function updateLiveButton() {
  try {
    const res = await fetch("/api/stream/status", { cache: "no-store" });
    const text = await res.text();
    const isLive = text.trim().toLowerCase() === "true";

    const btn = document.getElementById("live-btn");
    if (!btn) return;

    // normalize: remove any existing bg utilities first
    btn.classList.forEach(c => {
      if (c.startsWith("bg-") || c.startsWith("hover:bg-")) btn.classList.remove(c);
    });

    btn.classList.add(
      isLive ? "bg-red-600" : "bg-gray-800",
      isLive ? "hover:bg-red-500" : "hover:bg-gray-600",
      "text-gray-200","transition-colors"
    );
  } catch (e) {
    console.error("status fetch failed", e);
  }
}

document.addEventListener("DOMContentLoaded", () => {
  updateLiveButton();
  setInterval(updateLiveButton, 15000);
});
