// Supabase Auth Configuration - loaded from environment
let SUPABASE_URL;
let SUPABASE_ANON_KEY;

// Import Supabase from CDN
import { createClient } from 'https://cdn.skypack.dev/@supabase/supabase-js@2';

let supabase;

class SupabaseAuthManager {
  constructor() {
    this.currentUser = null;
    this.supabase = null;
    this.isSyncing = false;
    this.init();
  }

  async init() {
    try {
      // Fetch Supabase config from server
      const response = await fetch('/api/config/supabase');
      const config = await response.json();

      SUPABASE_URL = config.url;
      SUPABASE_ANON_KEY = config.anonKey;

      this.supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY);
      supabase = this.supabase; // Set global reference

      // Listen for auth changes
      this.supabase.auth.onAuthStateChange(async (event, session) => {
        console.log('Auth state changed:', event, session?.user?.email);
        this.currentUser = session?.user || null;
        this.updateUI(session?.user);

        if (session?.access_token) {
          this.setAuthToken(session.access_token);
        } else {
          this.clearAuthToken();
        }

        // Reconnect chat when auth state changes
        if (window.reconnectChat) {
          // Use a small delay to ensure token is properly set
          setTimeout(() => {
            console.log('🔄 Auth state changed, reconnecting chat...');
            window.reconnectChat();
          }, 100);
        }
      });

      // Get initial session AFTER setting up listener
      await this.getSession();

    } catch (error) {
      console.error('Failed to initialize Supabase:', error);
      this.showMessage('Failed to initialize authentication system', 'error');
    }
  }

  async getSession() {
    const { data: { session } } = await this.supabase.auth.getSession();
    if (session) {
      this.currentUser = session.user;
      this.setAuthToken(session.access_token);
      this.updateUI(session.user);
    }
  }

  async signUp(email, password) {
    try {
      const { data, error } = await this.supabase.auth.signUp({
        email: email,
        password: password,
      });

      if (error) throw error;

      return {
        success: true,
        user: data.user,
        message: data.user?.email_confirmed_at ? 'Account created!' : 'Check your email for verification!'
      };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async signIn(email, password) {
    try {
      const { data, error } = await this.supabase.auth.signInWithPassword({
        email: email,
        password: password,
      });

      if (error) throw error;

      return { success: true, user: data.user };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async signInWithTwitch() {
    try {
      // Use simple redirect to current origin + /auth/callback
      const baseUrl = window.location.origin;

      const { data, error } = await this.supabase.auth.signInWithOAuth({
        provider: 'twitch',
        options: {
          redirectTo: `${baseUrl}/auth/callback?returnTo=/live`
        }
      });

      if (error) throw error;

      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async signInWithGitHub() {
    try {
      // Use simple redirect to current origin + /auth/callback
      const baseUrl = window.location.origin;

      const { data, error } = await this.supabase.auth.signInWithOAuth({
        provider: 'github',
        options: {
          redirectTo: `${baseUrl}/auth/callback?returnTo=/live`
        }
      });

      if (error) throw error;

      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async signInWithGoogle() {
    try {
      // Use simple redirect to current origin + /auth/callback
      const baseUrl = window.location.origin;

      const { data, error } = await this.supabase.auth.signInWithOAuth({
        provider: 'google',
        options: {
          redirectTo: `${baseUrl}/auth/callback?returnTo=/live`
        }
      });

      if (error) throw error;

      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async signOut() {
    try {
      const { error } = await this.supabase.auth.signOut();
      if (error) throw error;

      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async resetPassword(email) {
    try {
      // Use current domain for redirect, but fallback to majesticcoding.com for production
      const baseUrl = window.location.hostname === 'localhost' || window.location.hostname.includes('127.0.0.1')
        ? window.location.origin
        : 'https://majesticcoding.com';

      const { data, error } = await this.supabase.auth.resetPasswordForEmail(email, {
        redirectTo: `${baseUrl}/auth/reset-password`
      });

      if (error) throw error;

      return { success: true, message: 'Password reset email sent!' };
    } catch (error) {
      return { success: false, error: error.message };
    }
  }

  async getAccessToken() {
    // First try to get the current session token from Supabase
    if (this.supabase) {
      const { data: { session } } = await this.supabase.auth.getSession();
      if (session?.access_token) {
        return session.access_token;
      }
    }

    // Fallback to localStorage
    return localStorage.getItem('supabase_token');
  }

  setAuthToken(token) {
    localStorage.setItem('supabase_token', token);
  }

  clearAuthToken() {
    localStorage.removeItem('supabase_token');
  }

  // Update UI based on auth state
  updateUI(user) {
    const loginSection = document.getElementById('login-section');
    const loggedInSection = document.getElementById('logged-in-section');
    const userEmail = document.getElementById('user-email');
    const userAvatar = document.getElementById('user-avatar');

    if (user) {
      // User is logged in
      if (loginSection) loginSection.style.display = 'none';
      if (loggedInSection) loggedInSection.style.display = 'block';
      if (userEmail) userEmail.textContent = user.email;
      if (userAvatar && user.user_metadata?.avatar_url) {
        userAvatar.src = user.user_metadata.avatar_url;
        userAvatar.style.display = 'block';
      }
    } else {
      // User is logged out
      if (loginSection) loginSection.style.display = 'block';
      if (loggedInSection) loggedInSection.style.display = 'none';
    }
  }

  // Helper to make authenticated API calls
  async apiCall(url, options = {}) {
    const token = await this.getAccessToken();

    if (token) {
      // Clean the token of any whitespace
      const cleanToken = token.trim();
      options.headers = {
        ...options.headers,
        'Authorization': `Bearer ${cleanToken}`
      };
    }

    return fetch(url, options);
  }

  // Get user info from API
  async getUserInfo() {
    try {
      const response = await this.apiCall('/api/user/info');

      if (response.ok) {
        return await response.json();
      } else {
        const error = await response.json();
        throw new Error(error.error || 'Failed to get user info');
      }
    } catch (error) {
      console.error('Get user info error:', error);
      throw error;
    }
  }

  // Sync user to database
  async syncUser() {
    try {
      const response = await this.apiCall('/api/user/sync', {
        method: 'POST'
      });

      if (response.ok) {
        return await response.json();
      } else {
        const error = await response.json();
        throw new Error(error.error || 'Failed to sync user');
      }
    } catch (error) {
      console.error('Sync user error:', error);
      throw error;
    }
  }

  // Get current user info
  getCurrentUser() {
    return this.currentUser;
  }

  // Helper to show messages
  showMessage(message, type = 'info') {
    const messageEl = document.getElementById('auth-message');
    if (messageEl) {
      messageEl.textContent = message;
      messageEl.className = `message ${type}`;
      messageEl.style.display = 'block';
      setTimeout(() => {
        messageEl.style.display = 'none';
      }, 5000);
    } else {
      console.log(`[${type.toUpperCase()}] ${message}`);
    }
  }
}

// Create global auth manager
window.authManager = new SupabaseAuthManager();

// Add event listeners when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  // Show/hide loading indicator
  function showLoading(show = true) {
    const buttons = document.querySelectorAll('button');
    buttons.forEach(btn => {
      btn.disabled = show;
      if (show) {
        btn.style.opacity = '0.6';
      } else {
        btn.style.opacity = '1';
      }
    });
  }

  // Show message to user
  function showMessage(message, type = 'info') {
    const messageEl = document.getElementById('auth-message');
    if (messageEl) {
      messageEl.textContent = message;
      messageEl.className = `message ${type}`;
      messageEl.style.display = 'block';
      setTimeout(() => {
        messageEl.style.display = 'none';
      }, 5000);
    } else {
      alert(message);
    }
  }

  // Sign in form
  const loginForm = document.getElementById('login-form');
  if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      showLoading(true);

      const email = document.getElementById('email').value;
      const password = document.getElementById('password').value;

      const result = await authManager.signIn(email, password);

      showLoading(false);

      if (result.success) {
        showMessage('Successfully signed in!', 'success');
      } else {
        showMessage('Sign in failed: ' + result.error, 'error');
      }
    });
  }

  // Sign up form
  const signupForm = document.getElementById('signup-form');
  if (signupForm) {
    signupForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      showLoading(true);

      const email = document.getElementById('signup-email').value;
      const password = document.getElementById('signup-password').value;

      const result = await authManager.signUp(email, password);

      showLoading(false);

      if (result.success) {
        showMessage(result.message, 'success');
      } else {
        showMessage('Sign up failed: ' + result.error, 'error');
      }
    });
  }

  // Twitch sign in
  const twitchLoginBtn = document.getElementById('twitch-login');
  if (twitchLoginBtn) {
    twitchLoginBtn.addEventListener('click', async () => {
      const result = await authManager.signInWithTwitch();
      if (!result.success) {
        showMessage('Twitch sign in failed: ' + result.error, 'error');
      }
    });
  }

  // GitHub sign in
  const githubLoginBtn = document.getElementById('github-login');
  if (githubLoginBtn) {
    githubLoginBtn.addEventListener('click', async () => {
      const result = await authManager.signInWithGitHub();
      if (!result.success) {
        showMessage('GitHub sign in failed: ' + result.error, 'error');
      }
    });
  }

  // Google sign in
  const googleLoginBtn = document.getElementById('google-login');
  if (googleLoginBtn) {
    googleLoginBtn.addEventListener('click', async () => {
      const result = await authManager.signInWithGoogle();
      if (!result.success) {
        showMessage('Google sign in failed: ' + result.error, 'error');
      }
    });
  }

  // Password reset
  const resetForm = document.getElementById('reset-form');
  if (resetForm) {
    resetForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      showLoading(true);

      const email = document.getElementById('reset-email').value;
      const result = await authManager.resetPassword(email);

      showLoading(false);

      if (result.success) {
        showMessage(result.message, 'success');
      } else {
        showMessage('Reset failed: ' + result.error, 'error');
      }
    });
  }

  // Sign out - handle multiple possible sign out buttons
  const signOutBtns = document.querySelectorAll('#sign-out, #sign-out-btn');
  signOutBtns.forEach(btn => {
    btn.addEventListener('click', async (e) => {
      e.preventDefault();
      console.log('Sign out clicked from:', btn.id);

      try {
        const result = await authManager.signOut();
        console.log('Sign out result:', result);

        if (result.success) {
          showMessage('Signed out successfully', 'success');
          // Redirect to home after a short delay
          setTimeout(() => {
            window.location.href = '/';
          }, 1000);
        } else {
          showMessage('Sign out failed: ' + (result.error || 'Unknown error'), 'error');
        }
      } catch (error) {
        console.error('Sign out error:', error);
        showMessage('Sign out error: ' + error.message, 'error');
      }
    });
  });
});

// Export for use in other scripts
window.supabase = supabase;