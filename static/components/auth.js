document.addEventListener('DOMContentLoaded', () => {
  function waitForClerk(callback) {
    if (window.Clerk) {
      callback();
    } else {
      setTimeout(() => waitForClerk(callback), 50);
    }
  }

  waitForClerk(async () => {
    await Clerk.load();

    const chatForm = document.getElementById('chat-form');
    if (chatForm) {
      chatForm.addEventListener('submit', function(e) {
        if (!Clerk.user) {
          e.preventDefault();
          window.location.href = "/auth";
        }
        // If signed in, allow submit as normal
      });
    }
  });
});