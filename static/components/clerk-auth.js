// Clerk Authentication Component

let clerkInstance = null;

function loadScript(src, attrs) {
  return new Promise(resolve => {
    const s = document.createElement('script');
    s.src = src;
    Object.entries(attrs || {}).forEach(([k, v]) => s.setAttribute(k, v));
    s.addEventListener('load', resolve);
    document.head.appendChild(s);
  });
}

async function cacheToken() {
  if (!window.Clerk?.session) {
    localStorage.removeItem('clerk_token');
    return;
  }
  try {
    const token = await window.Clerk.session.getToken();
    if (token) localStorage.setItem('clerk_token', token);
  } catch (_) {
    localStorage.removeItem('clerk_token');
  }
}

async function initClerk() {
  // Clear any stale token (e.g. leftover Supabase JWT) before loading Clerk.
  // cacheToken() below will repopulate it if the user has a valid Clerk session.
  localStorage.removeItem('clerk_token');

  const response = await fetch('/api/config/clerk');
  if (!response.ok) return;
  const { publishableKey } = await response.json();
  if (!publishableKey) return;

  const encoded = publishableKey.replace('pk_live_', '').replace('pk_test_', '');
  const clerkDomain = atob(encoded).replace(/\$$/, '');

  await loadScript(`https://${clerkDomain}/npm/@clerk/clerk-js@5/dist/clerk.browser.js`, {
    'data-clerk-publishable-key': publishableKey,
    crossorigin: 'anonymous',
  });

  await window.Clerk.load();
  clerkInstance = window.Clerk;

  await cacheToken();
  updateAuthUI();
  window.clerkAuthReady = true;
  document.dispatchEvent(new Event('clerk-auth-ready'));

  window.Clerk.addListener(async () => {
    await cacheToken();
    updateAuthUI();
    if (window.reconnectChat) setTimeout(window.reconnectChat, 100);
    document.dispatchEvent(new Event('clerk-auth-change'));
  });
}

function updateAuthUI() {
  const user = window.Clerk?.user;
  const redirectUrl = encodeURIComponent(window.location.pathname);

  document.querySelectorAll('[data-auth="login"]').forEach(el => {
    el.classList.toggle('hidden', !!user);
    if (el.tagName === 'A') el.href = `/auth?redirect_url=${redirectUrl}`;
  });
  document.querySelectorAll('[data-auth="logout"]').forEach(el =>
    el.classList.toggle('hidden', !user));
  document.querySelectorAll('[data-auth="email"]').forEach(el =>
    el.textContent = user?.primaryEmailAddress?.emailAddress || '');

  // Also patch any plain /auth links so they carry back the current page
  document.querySelectorAll('a[href="/auth"]').forEach(el =>
    el.setAttribute('href', `/auth?redirect_url=${redirectUrl}`));

  if (user) {
    localStorage.setItem('clerk_user_email', user.primaryEmailAddress?.emailAddress || '');
  } else {
    localStorage.removeItem('clerk_user_email');
    localStorage.removeItem('clerk_token');
  }
}

window.getAuthToken = async function () {
  await cacheToken();
  return localStorage.getItem('clerk_token');
};

window.getCurrentUser = function () {
  return window.Clerk?.user || null;
};

window.signOut = async function () {
  if (window.Clerk) await window.Clerk.signOut();
};

window.openSignIn = function () {
  if (window.Clerk) window.Clerk.redirectToSignIn({ redirectUrl: window.location.href });
};

initClerk();
