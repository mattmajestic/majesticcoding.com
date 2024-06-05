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

    if (Clerk.user) {
      const userButtonDiv = document.createElement("div");
      userButtonDiv.id = "user-button";
      if (document.getElementById("app")) {
        document.getElementById("app").appendChild(userButtonDiv);
      }
      Clerk.mountUserButton(userButtonDiv);
    } else {
      const signInDiv = document.createElement("div");
      signInDiv.id = "sign-in";
      if (document.getElementById("app")) {
        document.getElementById("app").appendChild(signInDiv);
      }
      Clerk.mountSignIn(signInDiv);
    }

    if (document.getElementById("header-content")) {
      if (Clerk.user) {
        document.getElementById("login-button").style.display = "none";
        document.getElementById("header-content").innerHTML += `
          <div id="user-button-header"></div>
        `;

        const userButtonHeaderDiv = document.getElementById("user-button-header");
        Clerk.mountUserButton(userButtonHeaderDiv);
      } else {
        document.getElementById("login-button").style.display = "block";
      }
    }
  });
};
