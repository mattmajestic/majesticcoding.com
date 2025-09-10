// Enhanced authentication handling with server-side session management

document.addEventListener('DOMContentLoaded', () => {
  let currentUser = null;

  function waitForClerk(callback) {
    if (window.Clerk) {
      callback();
    } else {
      setTimeout(() => waitForClerk(callback), 50);
    }
  }

  // Main authentication flow
  waitForClerk(async () => {
    await Clerk.load();

    const appDiv = document.getElementById("app");
    if (!appDiv) return;

    // Clear loading state
    appDiv.innerHTML = '';

    if (Clerk.user) {
      currentUser = Clerk.user;
      await handleAuthenticatedUser(appDiv);
    } else {
      await handleUnauthenticatedUser(appDiv);
    }

    // Listen for authentication state changes
    Clerk.addListener((resources) => {
      if (resources.user && !currentUser) {
        // User just signed in
        currentUser = resources.user;
        handleAuthenticatedUser(appDiv);
      } else if (!resources.user && currentUser) {
        // User just signed out
        currentUser = null;
        handleSignOut();
      }
    });
  });

  // Handle authenticated user - show user button and sync with server
  async function handleAuthenticatedUser(appDiv) {
    try {
      // Get JWT token for server communication
      const token = await Clerk.session.getToken();
      
      // Call our server login endpoint to establish server-side session
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      const authData = await response.json();
      
      if (authData.success) {
        // Show success message and user info
        showAuthenticatedState(appDiv, authData.user);
        
        // Redirect after a short delay
        setTimeout(() => {
          const urlParams = new URLSearchParams(window.location.search);
          const redirect = urlParams.get('redirect') || '/dashboard';
          window.location.href = redirect;
        }, 2000);
      } else {
        console.error('Server authentication failed:', authData.message);
        showUserButton(appDiv);
      }
    } catch (error) {
      console.error('Authentication error:', error);
      showUserButton(appDiv);
    }
  }

  // Handle unauthenticated user - show sign in form
  async function handleUnauthenticatedUser(appDiv) {
    const signInDiv = document.createElement("div");
    signInDiv.id = "sign-in";
    signInDiv.className = "w-full";
    appDiv.appendChild(signInDiv);

    const params = new URLSearchParams(window.location.search);
    const redirectPath = params.get("redirect") || "/dashboard";

    try {
      Clerk.mountSignIn(signInDiv, {
        redirectUrl: window.location.origin + "/auth?redirect=" + encodeURIComponent(redirectPath),
        afterSignInUrl: redirectPath,
        appearance: {
          elements: {
            formButtonPrimary: 'bg-blue-600 hover:bg-blue-700 text-sm font-medium',
            card: 'shadow-none border-0',
            headerTitle: 'text-xl font-semibold text-gray-800 dark:text-white',
            headerSubtitle: 'text-gray-600 dark:text-gray-300',
          }
        }
      });
    } catch (error) {
      console.error('Failed to mount sign in:', error);
      appDiv.innerHTML = `
        <div class="text-center text-red-600">
          <p>Error loading sign in form. Please refresh and try again.</p>
        </div>
      `;
    }
  }

  // Show authenticated state with user info
  function showAuthenticatedState(appDiv, userSession) {
    appDiv.innerHTML = `
      <div class="text-center">
        <div class="mb-6">
          <div class="w-16 h-16 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center mx-auto mb-4">
            <i class="fas fa-check-circle text-2xl text-green-600 dark:text-green-400"></i>
          </div>
          <h2 class="text-xl font-semibold text-gray-800 dark:text-white mb-2">Welcome back!</h2>
          <p class="text-gray-600 dark:text-gray-300">Successfully authenticated as</p>
          <p class="font-medium text-blue-600 dark:text-blue-400">${userSession.user.email}</p>
        </div>
        
        <div class="space-y-3">
          <div class="flex items-center justify-center text-sm text-gray-600 dark:text-gray-300">
            <i class="fas fa-spinner fa-spin mr-2"></i>
            Redirecting to dashboard...
          </div>
          
          <div id="user-button-container" class="flex justify-center"></div>
        </div>
      </div>
    `;

    // Mount user button in the container
    const userButtonContainer = document.getElementById('user-button-container');
    if (userButtonContainer) {
      Clerk.mountUserButton(userButtonContainer, {
        afterSignOutUrl: window.location.origin,
        appearance: {
          elements: {
            userButtonAvatarBox: 'w-8 h-8',
            userButtonPopoverCard: 'shadow-lg border border-gray-200 dark:border-gray-700',
          }
        }
      });
    }
  }

  // Show just the user button (fallback)
  function showUserButton(appDiv) {
    appDiv.innerHTML = '<div id="user-button-container" class="flex justify-center"></div>';
    
    const userButtonContainer = document.getElementById('user-button-container');
    if (userButtonContainer) {
      Clerk.mountUserButton(userButtonContainer, {
        afterSignOutUrl: window.location.origin,
        appearance: {
          elements: {
            userButtonAvatarBox: 'w-10 h-10',
          }
        }
      });
    }
  }

  // Handle sign out
  async function handleSignOut() {
    try {
      // Call server logout endpoint
      await fetch('/api/auth/logout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        }
      });
    } catch (error) {
      console.error('Logout error:', error);
    }

    // Redirect to home page
    window.location.href = '/';
  }

  // Utility function to get current user status
  window.getCurrentUser = async function() {
    if (!Clerk.user) return null;
    
    try {
      const token = await Clerk.session.getToken();
      const response = await fetch('/api/user/status', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      return await response.json();
    } catch (error) {
      console.error('Error getting user status:', error);
      return null;
    }
  };
});