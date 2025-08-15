document.addEventListener('DOMContentLoaded', async () => {
  if (!window.Clerk || !window.Clerk.CLERK_PUBLISHABLE_KEY) {
    console.error("Clerk or publishableKey missing");
    return;
  }

  await window.Clerk.load();

  window.Clerk.addListener(async ({ user }) => {
    if (!user) {
      console.error("No Clerk user found");
      return;
    }

    // Get Clerk session token
    const token = await window.Clerk.session.getToken();

    // Send token to Go backend
    fetch("/api/user/status", {
      method: "GET",
      headers: {
        "Authorization": `Bearer ${token}`
      }
    })
    .then(res => res.json())
    .then(data => {
      console.log("Backend response:", data);
      // Use backend user info as needed
    })
    .catch(err => console.error("Auth check failed:", err));
  });
});