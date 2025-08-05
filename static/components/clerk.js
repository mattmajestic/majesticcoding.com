// clerk.js

// Add Clerk initialization script dynamically
const clerkScript = document.createElement('script');
clerkScript.async = true;
clerkScript.crossOrigin = 'anonymous';
clerkScript.setAttribute('data-clerk-publishable-key', window.CLERK_PUBLISHABLE_KEY);
clerkScript.src = 'https://grand-falcon-32.clerk.accounts.dev/npm/@clerk/clerk-js@latest/dist/clerk.browser.js';
clerkScript.type = 'text/javascript';
document.head.appendChild(clerkScript);

clerkScript.onload = function() {
  window.addEventListener("load", async function () {
    await Clerk.load();

    if (Clerk.session) {
      const token = await Clerk.session.getToken();
      console.log("Clerk token:", token);
      fetch("/api/auth/status", {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${token}`
        }
      });
    }

    if (Clerk.user) {
      const userButtonDiv = document.createElement("div");
      userButtonDiv.id = "user-button";
      if (document.getElementById("app")) {
        document.getElementById("app").appendChild(userButtonDiv);
      }
      Clerk.mountUserButton(userButtonDiv, {
        afterSignOutUrl: "http://localhost:8080"
      });
    } else {
      const signInDiv = document.createElement("div");
      signInDiv.id = "sign-in";

      // Get ?redirect= from URL or default to "/"
      const params = new URLSearchParams(window.location.search);
      const redirectPath = params.get("redirect") || "/";

      Clerk.mountSignIn(signInDiv, {
        redirectUrl: redirectPath,
        afterSignInUrl: redirectPath
      });
    }

    if (document.getElementById("header-content")) {
      if (Clerk.user) {
        document.getElementById("login-button").style.display = "none";
        document.getElementById("header-content").innerHTML += `
          <div id="user-button-header"></div>
        `;

        const userButtonHeaderDiv = document.getElementById("user-button-header");
        Clerk.mountUserButton(userButtonHeaderDiv, {
          afterSignOutUrl: "http://localhost:8080"
        });
      } else {
        document.getElementById("login-button").style.display = "block";
      }
    }
  });
};
