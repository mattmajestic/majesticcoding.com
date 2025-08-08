document.addEventListener('DOMContentLoaded', async () => {
  if (!window.Clerk || !window.Clerk.CLERK_PUBLISHABLE_KEY) {
    console.error("Clerk or publishableKey missing");
    return;
  }

  await window.Clerk.load();

  // Wait for Clerk to be ready
  window.Clerk.addListener(({ user }) => {
    if (!user) {
      console.error("No Clerk user found");
      return;
    }

    const sessionData = {
      id: user.id,
      email: user.primaryEmailAddress?.emailAddress || "",
      username: user.username || user.firstName || "anon",
    };

    // Use sessionData as needed
    console.log("User:", sessionData);
  });
});