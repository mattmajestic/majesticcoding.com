// Clerk Authentication for Login & Out

// Dynamically load Clerk script
const clerkScript = document.createElement('script');
clerkScript.async = true;
clerkScript.crossOrigin = 'anonymous';
clerkScript.setAttribute('data-clerk-publishable-key', window.CLERK_PUBLISHABLE_KEY);
clerkScript.src = 'https://grand-falcon-32.clerk.accounts.dev/npm/@clerk/clerk-js@latest/dist/clerk.browser.js';
clerkScript.type = 'text/javascript';
document.head.appendChild(clerkScript);

// On Clerk script load
clerkScript.onload = function () {
  window.addEventListener("load", async function () {
    await Clerk.load();

    const appDiv = document.getElementById("app");
    if (!appDiv) return;

    if (Clerk.user) {
      const userButtonDiv = document.createElement("div");
      userButtonDiv.id = "user-button";
      appDiv.appendChild(userButtonDiv);

      Clerk.mountUserButton(userButtonDiv, {
        afterSignOutUrl: "http://localhost:8080"
      });
    } else {
      const signInDiv = document.createElement("div");
      signInDiv.id = "sign-in";
      appDiv.appendChild(signInDiv);

      const params = new URLSearchParams(window.location.search);
      const redirectPath = params.get("redirect") || "/";

      Clerk.mountSignIn(signInDiv, {
        redirectUrl: redirectPath,
        afterSignInUrl: redirectPath
      });
    }
  });
};
