document.addEventListener("DOMContentLoaded", () => {
  const buttons = document.querySelectorAll(".stats-button");
  const loader = document.getElementById("stats-loader");

  buttons.forEach(button => {
    button.addEventListener("click", async () => {
      const provider = button.getAttribute("data-provider");
      if (!provider || !loader) return;

      loader.classList.remove("hidden");

      try {
        const res = await fetch(`/api/stats/${provider}`);
        const data = await res.json();
        loader.classList.add("hidden");
        openModal(provider, data);
      } catch (err) {
        loader.classList.add("hidden");
        openModal("Error", { message: "Failed to load stats." });
      }
    });
  });
});

function openModal(title, data) {
  document.getElementById("modal-title").textContent =
    `${title.charAt(0).toUpperCase() + title.slice(1)} Stats via API`;

  const content = Object.entries(data)
    .map(([k, v]) => {
      let formattedValue = v;

      // Try formatting if it's a number or a string that looks like a number
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
  document.getElementById("stats-modal").classList.remove("hidden");
}


function closeModal() {
  document.getElementById("stats-modal").classList.add("hidden");
}