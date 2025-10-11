// Stats Modal from API Data 
// of YouTube, Twitch Github & Leetcode

document.addEventListener("DOMContentLoaded", () => {
  // Prevent double initialization
  if (window.statsInitialized) return;
  window.statsInitialized = true;

  const buttons = document.querySelectorAll(".stats-button");

  buttons.forEach(button => {
    button.addEventListener("click", async () => {
      const provider = button.getAttribute("data-provider");
      if (!provider) return;

      showModalShell(provider); // Show modal immediately with "Loading..."

      try {
        const res = await fetch(`/api/stats/${provider}`);
        const data = await res.json();
        populateModalContent(data); // Fill in once loaded
      } catch (err) {
        populateModalContent({ error: "Failed to load stats." });
      }
    });
  });
});

function showModalShell(title) {
  document.getElementById("modal-title").textContent =
    `${title.charAt(0).toUpperCase() + title.slice(1)} Stats via API`;

  document.getElementById("stats-content").innerHTML = `
    <div class="flex items-center justify-center py-8 gap-3">
      <i class="fas fa-spinner fa-spin text-3xl text-green-400"></i>
      <span class="text-xl text-green-300 font-semibold">Querying API...</span>
    </div>
  `;

  document.getElementById("stats-modal").classList.remove("hidden");
  
  // Add click outside to close functionality
  setTimeout(() => {
    document.addEventListener("click", handleOutsideClick);
  }, 100); // Small delay to prevent immediate closing
}

function populateModalContent(data) {
  const content = Object.entries(data)
    .map(([k, v]) => {
      let formattedValue = v;
      if (!isNaN(v)) {
        formattedValue = Number(v).toLocaleString();
      }

      // Format camelCase labels to have spaces (e.g., MainLanguages -> Main Languages)
      const formattedKey = k.replace(/([A-Z])/g, ' $1').trim();

      return `
        <div class="stat-box glowing-effect">
          <div class="stat-key">${formattedKey}</div>
          <div class="stat-value">${formattedValue}</div>
        </div>
      `;
    })
    .join("");

  document.getElementById("stats-content").innerHTML = content;
}

function closeModal() {
  document.getElementById("stats-modal").classList.add("hidden");
  // Remove the outside click listener when modal is closed
  document.removeEventListener("click", handleOutsideClick);
}

function handleOutsideClick(event) {
  const modal = document.getElementById("stats-modal");
  const modalContent = modal.querySelector("div");
  
  // Check if modal is visible and click is outside modal content
  if (!modal.classList.contains("hidden") && modalContent && !modalContent.contains(event.target)) {
    closeModal();
  }
}
