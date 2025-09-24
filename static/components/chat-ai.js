// AI Chat Interface
class AIChatInterface {
  constructor() {
    this.messageCount = 0;
    this.isLoading = false;
    this.conversationHistory = [];

    this.elements = {
      chatMessages: document.getElementById('chat-messages'),
      aiInput: document.getElementById('ai-input'),
      aiForm: document.getElementById('ai-chat-form'),
      sendButton: document.getElementById('send-button'),
      clearButton: document.getElementById('clear-chat'),
      authStatus: document.getElementById('auth-status'),
      providerSelect: document.getElementById('ai-provider'),
      messageCountEl: document.getElementById('message-count'),
      currentProviderEl: document.getElementById('current-provider')
    };

    this.initializeEventListeners();
    this.checkAuthStatus();
    this.loadProviders();
    this.setupQuickActions();
    this.setupTextarea();
  }

  initializeEventListeners() {
    // Form submission
    this.elements.aiForm.addEventListener('submit', (e) => {
      e.preventDefault();
      this.sendMessage();
    });

    // Clear chat
    this.elements.clearButton.addEventListener('click', () => {
      this.clearChat();
    });

    // Provider selection
    this.elements.providerSelect.addEventListener('change', () => {
      this.updateCurrentProvider();
    });

    // Enter to send, Shift+Enter for new line
    this.elements.aiInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        if (!this.isLoading && this.elements.aiInput.value.trim()) {
          this.sendMessage();
        }
      }
    });

    // Input validation
    this.elements.aiInput.addEventListener('input', () => {
      this.validateInput();
      this.updateCharCounter();
    });
  }

  setupTextarea() {
    // Auto-resize textarea
    this.elements.aiInput.addEventListener('input', () => {
      const textarea = this.elements.aiInput;
      textarea.style.height = '48px';
      textarea.style.height = Math.min(textarea.scrollHeight, 120) + 'px';
    });
  }

  setupQuickActions() {
    const quickActions = document.querySelectorAll('.quick-action');
    quickActions.forEach(button => {
      button.addEventListener('click', () => {
        const prompt = button.getAttribute('data-prompt');
        this.elements.aiInput.value = prompt;
        this.validateInput();
        this.elements.aiInput.focus();
      });
    });
  }

  updateCharCounter() {
    const length = this.elements.aiInput.value.length;
    this.elements.charCounter.textContent = length > 0 ? `${length}` : '';
  }

  validateInput() {
    const hasText = this.elements.aiInput.value.trim().length > 0;
    const hasAuth = this.getAuthToken();

    this.elements.sendButton.disabled = !hasText || !hasAuth || this.isLoading;

    if (!hasAuth) {
      this.elements.sendButton.title = 'Please log in to use AI chat';
    } else if (!hasText) {
      this.elements.sendButton.title = 'Enter a message';
    } else if (this.isLoading) {
      this.elements.sendButton.title = 'Please wait...';
    } else {
      this.elements.sendButton.title = 'Send message';
    }
  }

  async checkAuthStatus() {
    const token = this.getAuthToken();

    if (token) {
      this.elements.authStatus.textContent = '✅ Authenticated';
      this.elements.authStatus.className = 'px-3 py-1 rounded text-xs bg-green-900 text-green-300 hover:bg-green-800 transition-colors cursor-pointer';
      this.elements.authStatus.href = '/settings';
      this.elements.authStatus.title = 'View Profile & Settings';
    } else {
      this.elements.authStatus.textContent = 'Sign In';
      this.elements.authStatus.className = 'px-3 py-1 rounded text-xs bg-red-900 text-red-300 hover:bg-red-800 transition-colors cursor-pointer';
      this.elements.authStatus.href = '/auth';
      this.elements.authStatus.title = 'Click to sign in';
      this.showAuthMessage();
    }

    this.validateInput();
  }

  getAuthToken() {
    // Try localStorage first
    let token = localStorage.getItem('supabase_token');

    // Fallback to Supabase session if available
    if (!token && window.supabaseAuth) {
      token = window.supabaseAuth.getAuthToken();
    }

    return token;
  }

  showAuthMessage() {
    const authMessage = this.createMessage('system', 'Please log in to start chatting with AI. Click "Sign In" in the top right corner.', 'System');
    this.elements.chatMessages.appendChild(authMessage);
    this.scrollToBottom();
  }

  async loadProviders() {
    try {
      const token = this.getAuthToken();
      if (!token) return;

      const response = await fetch('/api/llm/providers', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        const data = await response.json();
        this.updateProviderOptions(data.providers, data.fallback);
      }
    } catch (error) {
      console.error('Failed to load providers:', error);
    }
  }

  updateProviderOptions(providers, fallback) {
    const select = this.elements.providerSelect;

    // Clear existing options except "Auto"
    while (select.children.length > 1) {
      select.removeChild(select.lastChild);
    }

    // Add available providers
    providers.forEach(provider => {
      const option = document.createElement('option');
      option.value = provider;
      option.textContent = this.formatProviderName(provider);
      select.appendChild(option);
    });

    // Update current provider display
    this.elements.currentProviderEl.textContent = fallback ? this.formatProviderName(fallback) : 'Auto';
  }

  formatProviderName(provider) {
    const names = {
      'gemini': 'Google Gemini',
      'anthropic': 'Anthropic Claude',
      'openai': 'OpenAI GPT',
      'groq': 'Groq Llama'
    };
    return names[provider] || provider;
  }

  updateCurrentProvider() {
    const selected = this.elements.providerSelect.value;
    this.elements.currentProviderEl.textContent = selected ? this.formatProviderName(selected) : 'Auto';
  }

  async sendMessage() {
    const message = this.elements.aiInput.value.trim();
    if (!message || this.isLoading) return;

    const token = this.getAuthToken();
    if (!token) {
      this.checkAuthStatus();
      return;
    }

    // Add user message
    const userMessage = this.createMessage('user', message, 'You');
    this.elements.chatMessages.appendChild(userMessage);

    // Clear input and show loading
    this.elements.aiInput.value = '';
    this.elements.aiInput.style.height = '48px';
    this.setLoading(true);
    this.scrollToBottom();

    try {
      const provider = this.elements.providerSelect.value || undefined;

      const response = await fetch('/api/llm/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          prompt: message,
          provider: provider
        })
      });

      const data = await response.json();

      if (response.ok) {
        // Add AI response
        const aiMessage = this.createMessage('assistant', data.response, `AI (${data.provider})`);
        this.elements.chatMessages.appendChild(aiMessage);

        // Update current provider
        this.elements.currentProviderEl.textContent = this.formatProviderName(data.provider);
      } else {
        throw new Error(data.error || 'Failed to get AI response');
      }

    } catch (error) {
      console.error('AI chat error:', error);
      const errorMessage = this.createMessage('error', `Error: ${error.message}`, 'System');
      this.elements.chatMessages.appendChild(errorMessage);
    } finally {
      this.setLoading(false);
      this.messageCount++;
      this.elements.messageCountEl.textContent = this.messageCount;
      this.scrollToBottom();
      this.elements.aiInput.focus();
    }
  }

  createMessage(type, content, author) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message-${type} flex ${type === 'user' ? 'justify-end' : 'justify-start'} mb-4`;

    const messageContainer = document.createElement('div');
    messageContainer.className = `flex items-start space-x-3 max-w-[80%] ${type === 'user' ? 'flex-row-reverse space-x-reverse' : ''}`;

    // Avatar
    const avatar = document.createElement('img');
    if (type === 'assistant') {
      avatar.src = 'https://avatars.githubusercontent.com/u/33904170?v=4';
      avatar.alt = 'AI Assistant';
      avatar.className = 'w-8 h-8 rounded-full flex-shrink-0 mt-1';
    } else {
      avatar.src = 'https://via.placeholder.com/32/75a7da/FFFFFF?text=U';
      avatar.alt = 'User';
      avatar.className = 'w-8 h-8 rounded-full flex-shrink-0 mt-1 border-2 border-white shadow-md';
    }

    const bubble = document.createElement('div');
    bubble.className = this.getMessageStyles(type);

    const header = document.createElement('div');
    header.className = 'flex items-center justify-between mb-2 text-xs opacity-75';

    const authorSpan = document.createElement('span');
    authorSpan.className = 'font-medium';
    authorSpan.textContent = author;

    const timeSpan = document.createElement('span');
    timeSpan.className = 'opacity-60';
    timeSpan.textContent = new Date().toLocaleTimeString();

    header.appendChild(authorSpan);
    header.appendChild(timeSpan);

    const contentDiv = document.createElement('div');
    contentDiv.className = 'text-sm leading-relaxed';

    // Always apply formatting (renamed to formatText since it handles more than just code)
    contentDiv.innerHTML = this.formatText(content);

    bubble.appendChild(header);
    bubble.appendChild(contentDiv);

    messageContainer.appendChild(avatar);
    messageContainer.appendChild(bubble);
    messageDiv.appendChild(messageContainer);

    return messageDiv;
  }

  getMessageStyles(type) {
    switch (type) {
      case 'user':
        return 'message-user-bubble';
      case 'assistant':
        return 'message-assistant-bubble';
      case 'system':
        return 'message-system-bubble';
      case 'error':
        return 'message-error-bubble';
      default:
        return 'message-assistant-bubble';
    }
  }

  formatText(text) {
    return text
      // Code blocks first (multi-line) - before other processing
      .replace(/```(\w+)?\n([\s\S]*?)```/g, '<pre class="ai-code-block">$2</pre>')
      // Inline code - before other processing to avoid conflicts
      .replace(/`([^`]+)`/g, '<code class="ai-inline-code">$1</code>')
      // Bold text
      .replace(/\*\*([^*]+)\*\*/g, '<strong class="font-bold">$1</strong>')
      // Bulleted lists (lines starting with * or -, but not bold markers)
      .replace(/^[\*\-]\s+(.+)$/gm, '<div class="flex items-start gap-2 my-1"><span class="text-light-blue mt-1">•</span><span>$1</span></div>')
      // Numbered lists (lines starting with numbers)
      .replace(/^(\d+)\.?\s+(.+)$/gm, '<div class="flex items-start gap-2 my-1"><span class="text-light-blue font-medium min-w-6">$1.</span><span>$2</span></div>')
      // Headers (lines starting with #)
      .replace(/^#{1}\s+(.+)$/gm, '<h1 class="text-lg font-bold mt-4 mb-2 text-white">$1</h1>')
      .replace(/^#{2}\s+(.+)$/gm, '<h2 class="text-base font-bold mt-3 mb-1 text-white">$1</h2>')
      .replace(/^#{3}\s+(.+)$/gm, '<h3 class="text-sm font-bold mt-2 mb-1 text-white">$1</h3>')
      // Italic text (after bold to avoid conflicts)
      .replace(/\*([^*]+)\*/g, '<em class="italic">$1</em>')
      // Line breaks last
      .replace(/\n/g, '<br>');
  }

  setLoading(loading) {
    this.isLoading = loading;

    if (loading) {
      // Add thinking message to chat
      this.thinkingMessage = this.createThinkingMessage();
      this.elements.chatMessages.appendChild(this.thinkingMessage);
      this.scrollToBottom();

      this.elements.sendButton.disabled = true;
    } else {
      // Remove thinking message
      if (this.thinkingMessage) {
        this.thinkingMessage.remove();
        this.thinkingMessage = null;
      }
    }

    this.validateInput();
  }

  createThinkingMessage() {
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message-assistant flex justify-start mb-4';

    const messageContainer = document.createElement('div');
    messageContainer.className = 'flex items-start space-x-3 max-w-[80%]';

    // AI Avatar
    const avatar = document.createElement('img');
    avatar.src = 'https://avatars.githubusercontent.com/u/33904170?v=4';
    avatar.alt = 'AI Assistant';
    avatar.className = 'w-8 h-8 rounded-full flex-shrink-0 mt-1';

    const bubble = document.createElement('div');
    bubble.className = 'message-assistant-bubble flex items-center space-x-2';

    // Animated thinking dots
    const thinkingDiv = document.createElement('div');
    thinkingDiv.innerHTML = `
      <div class="flex items-center space-x-1">
        <span>AI is thinking</span>
        <div class="flex space-x-1">
          <div class="w-2 h-2 bg-current rounded-full animate-bounce" style="animation-delay: 0ms"></div>
          <div class="w-2 h-2 bg-current rounded-full animate-bounce" style="animation-delay: 150ms"></div>
          <div class="w-2 h-2 bg-current rounded-full animate-bounce" style="animation-delay: 300ms"></div>
        </div>
      </div>
    `;

    bubble.appendChild(thinkingDiv);
    messageContainer.appendChild(avatar);
    messageContainer.appendChild(bubble);
    messageDiv.appendChild(messageContainer);

    return messageDiv;
  }

  clearChat() {
    // Remove all messages except welcome
    const messages = this.elements.chatMessages.querySelectorAll('.message-user, .message-assistant, .message-error, .message-system');
    messages.forEach(msg => msg.remove());

    this.messageCount = 0;
    this.elements.messageCountEl.textContent = '0';
    this.conversationHistory = [];
  }

  scrollToBottom() {
    // Use setTimeout to ensure DOM has updated
    setTimeout(() => {
      this.elements.chatMessages.scrollTo({
        top: this.elements.chatMessages.scrollHeight,
        behavior: 'smooth'
      });
    }, 100);
  }
}

// Initialize AI Chat Interface when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  window.aiChat = new AIChatInterface();
});

// Re-check auth when Supabase auth changes
document.addEventListener('supabase-auth-change', () => {
  if (window.aiChat) {
    window.aiChat.checkAuthStatus();
    window.aiChat.loadProviders();
  }
});