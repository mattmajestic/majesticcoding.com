// Stats Modal from API Data 
// of YouTube, Twitch Github & Leetcode

document.addEventListener("DOMContentLoaded", () => {
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
}

function populateModalContent(data) {
  const content = Object.entries(data)
    .map(([k, v]) => {
      let formattedValue = v;
      if (!isNaN(v)) {
        formattedValue = Number(v).toLocaleString();
      }
      return `
        <div class="stat-box glowing-effect">
          <div class="stat-key">${k}</div>
          <div class="stat-value">${formattedValue}</div>
        </div>
      `;
    })
    .join("");

  document.getElementById("stats-content").innerHTML = content;
}

function closeModal() {
  document.getElementById("stats-modal").classList.add("hidden");
}
