export async function login(containerId) {
  try {
    await Clerk.load();
    const container = document.getElementById(containerId);

    if (!container) {
      console.error(`❌ Element with ID '${containerId}' not found.`);
      return;
    }

    Clerk.mountSignIn(container);
  } catch (err) {
    console.error("❌ Clerk login failed:", err);
  }
}

export function logout(redirectUrl = "/login") {
  Clerk.signOut({ redirectUrl });
}