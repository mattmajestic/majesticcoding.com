import { clearElement } from "../utils/dom.js";

export async function displayUser(containerId) {
  try {
    await Clerk.load();
    const user = await Clerk.user;
    const container = document.getElementById(containerId);

    if (!container) {
      console.error(`❌ Element with ID '${containerId}' not found.`);
      return;
    }

    clearElement(container);

    if (user) {
      const avatar = document.createElement("img");
      avatar.src = user.imageUrl;
      avatar.alt = "User Avatar";
      avatar.className = "rounded-full w-16 h-16 mx-auto mb-2";

      const logoutBtn = document.createElement("button");
      logoutBtn.textContent = "Logout";
      logoutBtn.className = "bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded";
      logoutBtn.onclick = () => Clerk.signOut({ redirectUrl: "/login" });

      container.appendChild(avatar);
      container.appendChild(logoutBtn);
    } else {
      container.innerHTML = "<p class='text-red-500'>User not logged in.</p>";
    }
  } catch (err) {
    console.error("❌ Failed to display user:", err);
  }
}