(function () {
  const overlay = document.getElementById("stream-chat");
  overlay.style.display = "flex";
  overlay.style.flexDirection = "row";
  overlay.style.gap = "16px";
  overlay.style.alignItems = "center";
  overlay.style.justifyContent = "center";

  const users = {};

  // Function to add or update a user
  function addUser(username, avatar) {
    if (!users[username]) {
      // Create user element
      const userDiv = document.createElement("div");
      userDiv.style.display = "flex";
      userDiv.style.flexDirection = "column";
      userDiv.style.alignItems = "center";
      userDiv.style.gap = "8px";

      const img = document.createElement("img");
      img.src = avatar || "/static/img/default-avatar.png"; // Default avatar
      img.style.width = "60px";
      img.style.height = "60px";
      img.style.borderRadius = "50%";
      img.style.border = "2px solid #fff";
      img.style.boxShadow = "0 2px 6px rgba(0, 0, 0, 0.5)";
      userDiv.appendChild(img);

      const name = document.createElement("div");
      name.textContent = username;
      name.style.color = "#fff";
      name.style.fontSize = "14px";
      name.style.fontWeight = "bold";
      userDiv.appendChild(name);

      overlay.appendChild(userDiv);
      users[username] = userDiv;

      // Remove user after 30 seconds of inactivity
      setTimeout(() => {
        if (users[username]) {
          overlay.removeChild(users[username]);
          delete users[username];
        }
      }, 30000);
    }
  }

  // WebSocket connection to Twitch messages
  const ws = new WebSocket(
    (location.protocol === "https:" ? "wss" : "ws") + "://" + location.host + "/ws/twitch"
  );

  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data);
      const username = msg.display_name || msg.username || "Unknown";
      const avatar = msg.avatar || `/static/img/default-avatar.png`;
      addUser(username, avatar);
    } catch (err) {
      console.error("Failed to parse WebSocket message:", err);
    }
  };

  ws.onerror = (err) => {
    console.error("WebSocket error:", err);
  };
})();