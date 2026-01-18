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
      micButton: document.getElementById('mic-button'),
      micMenuButton: document.getElementById('mic-menu-button'),
      micDeviceMenu: document.getElementById('mic-device-menu'),
      clearButton: document.getElementById('clear-chat'),
      authStatus: document.getElementById('auth-status'),
      providerSelect: document.getElementById('ai-provider'),
      messageCountEl: document.getElementById('message-count'),
      currentProviderEl: document.getElementById('current-provider'),
      charCounter: document.getElementById('char-counter') // Add missing element
    };

    this.speechState = {
      supported: false,
      recognition: null,
      isListening: false,
      baseText: '',
      isStarting: false,
      suppressErrors: false,
      isInterim: false
    };

    this.mediaState = {
      supported: false,
      recorder: null,
      isRecording: false,
      isStreaming: false,
      chunks: [],
      stream: null,
      mimeType: '',
      baseText: '',
      deviceId: localStorage.getItem('ai_mic_device_id') || '',
      streamSocket: null,
      pendingChunks: [],
      lastStreamText: '',
      audioContext: null,
      analyser: null,
      monitorId: null,
      hasSpoken: false,
      lastSpeechAt: 0
    };

    this.initializeEventListeners();
    this.setupSpeechRecognition();
    this.setupMediaRecorder();
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

    // Voice input
    if (this.elements.micButton) {
      this.elements.micButton.addEventListener('click', () => {
        this.toggleVoiceInput();
      });
    }

    if (this.elements.micMenuButton) {
      this.elements.micMenuButton.addEventListener('click', async () => {
        await this.ensureMicDevices();
        this.toggleMicDeviceSelect();
      });
    }

    if (this.elements.micDeviceMenu) {
      document.addEventListener('click', (event) => {
        if (!this.elements.micDeviceMenu) return;
        if (!this.elements.micMenuButton) return;
        const menu = this.elements.micDeviceMenu;
        const button = this.elements.micMenuButton;
        const micButton = this.elements.micButton;

        if (menu.classList.contains('hidden')) return;
        if (menu.contains(event.target) || button.contains(event.target) || (micButton && micButton.contains(event.target))) {
          return;
        }
        this.toggleMicDeviceSelect(false);
      });
    }
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
    if (this.elements.charCounter) {
      const length = this.elements.aiInput.value.length;
      this.elements.charCounter.textContent = length > 0 ? `${length}` : '';
    }
  }

  setupSpeechRecognition() {
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    if (!SpeechRecognition || !this.elements.micButton) return;

    const recognition = new SpeechRecognition();
    recognition.continuous = true;
    recognition.interimResults = true;
    recognition.lang = 'en-US';

    recognition.addEventListener('result', (event) => {
      let interimText = '';
      let finalText = '';

      for (let i = event.resultIndex; i < event.results.length; i++) {
        const transcript = event.results[i][0].transcript;
        if (event.results[i].isFinal) {
          finalText += transcript;
        } else {
          interimText += transcript;
        }
      }

      const base = this.speechState.baseText.trim();
      const combined = [base, finalText + interimText].filter(Boolean).join(' ').trim();
      this.elements.aiInput.value = combined;
      this.validateInput();
      this.updateCharCounter();
    });

    recognition.addEventListener('start', () => {
      if (this.speechState.isListening) {
        this.setVoiceListening(true);
      }
    });

    recognition.addEventListener('speechstart', () => {
      this.setVoiceSpeaking(true);
    });

    recognition.addEventListener('speechend', () => {
      this.setVoiceSpeaking(false);
    });

    recognition.addEventListener('end', () => {
      if (!this.speechState.isListening) return;

      // Chrome ends on silence; auto-restart for continuous listening.
      setTimeout(() => {
        if (!this.speechState.isListening) return;
        this.safeStartRecognition();
      }, 150);
    });

    recognition.addEventListener('error', (event) => {
      console.error('Speech recognition error:', event.error);
      if (this.speechState.suppressErrors) {
        this.speechState.isListening = false;
        this.speechState.isInterim = false;
        this.speechState.suppressErrors = false;
        return;
      }
      if (event.error === 'not-allowed' || event.error === 'service-not-allowed') {
        this.speechState.isListening = false;
        this.setVoiceListening(false);
        return;
      }

      if (event.error === 'network') {
        this.speechState.isListening = false;
        this.setVoiceListening(false);
        this.showSpeechError('Speech service unavailable. Check your connection and try again.');
        return;
      }

      if (this.speechState.isListening) {
        setTimeout(() => {
          if (!this.speechState.isListening) return;
          this.safeStartRecognition();
        }, 200);
      }
    });

    this.speechState.supported = true;
    this.speechState.recognition = recognition;
  }

  setupMediaRecorder() {
    if (!this.elements.micButton) return;
    if (!navigator.mediaDevices || !window.MediaRecorder) return;

    const preferredTypes = [
      'audio/webm;codecs=opus',
      'audio/webm',
      'audio/ogg;codecs=opus',
      'audio/ogg'
    ];

    for (const type of preferredTypes) {
      if (MediaRecorder.isTypeSupported(type)) {
        this.mediaState.mimeType = type;
        break;
      }
    }

    this.mediaState.supported = true;
  }

  getSTTMode() {
    const forcedMode = window.AI_CHAT_STT;
    if (forcedMode === 'browser' || forcedMode === 'server' || forcedMode === 'server_stream') return forcedMode;
    if (this.speechState.supported) return 'browser';
    if (this.mediaState.supported) return 'server_stream';
    return 'none';
  }

  async ensureMicDevices() {
    if (!this.elements.micDeviceMenu) return;
    if (!navigator.mediaDevices || !navigator.mediaDevices.enumerateDevices) return;

    const existingOptions = this.elements.micDeviceMenu.querySelectorAll('[data-device-id]');
    if (existingOptions.length > 0) return;

    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      stream.getTracks().forEach(track => track.stop());
    } catch (error) {
      console.error('Microphone permission error:', error);
      return;
    }

    const devices = await navigator.mediaDevices.enumerateDevices();
    const audioInputs = devices.filter(device => device.kind === 'audioinput');

    if (audioInputs.length === 0) return;

    this.elements.micDeviceMenu.innerHTML = `
      <div class="px-3 py-2 text-xs text-gray-200 bg-gray-700">Microphones</div>
    `;

    audioInputs.forEach((device, index) => {
      const button = document.createElement('button');
      button.type = 'button';
      button.dataset.deviceId = device.deviceId;
      button.className = 'w-full text-left px-3 py-3 text-sm text-gray-100 bg-gray-700 hover:bg-gray-600 transition-colors flex items-center justify-between';
      const label = device.label || `Microphone ${index + 1}`;
      button.innerHTML = `<span class="truncate">${label}</span><span class="text-green-400 text-xs hidden" data-selected>✓</span>`;
      button.addEventListener('click', () => {
        this.mediaState.deviceId = device.deviceId;
        localStorage.setItem('ai_mic_device_id', device.deviceId);
        this.setMicMenuTitle(label);
        this.updateMicMenuSelection(device.deviceId);
        this.toggleMicDeviceSelect(false);
      });
      this.elements.micDeviceMenu.appendChild(button);
    });

    let selected = audioInputs.find(device => device.deviceId === this.mediaState.deviceId);
    if (!selected && audioInputs.length > 0) {
      selected = audioInputs[0];
      this.mediaState.deviceId = selected.deviceId;
      localStorage.setItem('ai_mic_device_id', selected.deviceId);
    }
    if (selected) {
      this.setMicMenuTitle(selected.label || 'Microphone');
      this.updateMicMenuSelection(selected.deviceId);
    }
  }

  toggleMicDeviceSelect(forceState) {
    if (!this.elements.micDeviceMenu) return;
    if (typeof forceState === 'boolean') {
      this.elements.micDeviceMenu.classList.toggle('hidden', !forceState);
      return;
    }
    this.elements.micDeviceMenu.classList.toggle('hidden');
  }

  setMicMenuTitle(label) {
    if (!this.elements.micMenuButton) return;
    this.elements.micMenuButton.title = `Microphone: ${label}`;
  }

  updateMicMenuSelection(deviceId) {
    if (!this.elements.micDeviceMenu) return;
    const items = this.elements.micDeviceMenu.querySelectorAll('[data-device-id]');
    items.forEach(item => {
      const selectedBadge = item.querySelector('[data-selected]');
      if (selectedBadge) {
        selectedBadge.classList.toggle('hidden', item.dataset.deviceId !== deviceId);
      }
      item.classList.toggle('bg-gray-600', item.dataset.deviceId === deviceId);
    });
  }

  toggleVoiceInput() {
    const mode = this.getSTTMode();
    if (mode === 'browser') {
      this.toggleBrowserSpeech();
      return;
    }

    if (mode === 'server') {
      this.toggleServerRecording();
      return;
    }

    if (mode === 'server_stream') {
      this.toggleServerStreamingRecording();
      return;
    }

    if (this.elements.micButton) {
      this.elements.micButton.disabled = true;
      this.elements.micButton.title = 'Voice input not supported in this browser';
      this.elements.micButton.classList.add('opacity-50', 'cursor-not-allowed');
    }
  }

  toggleBrowserSpeech() {
    if (!this.speechState.supported || !this.speechState.recognition) return;

    if (this.speechState.isListening) {
      this.speechState.recognition.stop();
      this.speechState.isListening = false;
      this.setVoiceListening(false);
      return;
    }

    this.speechState.baseText = this.elements.aiInput.value.trim();
    this.speechState.isListening = true;
    this.speechState.isInterim = false;
    this.speechState.suppressErrors = false;
    this.setVoiceListening(true);
    this.safeStartRecognition();
    this.elements.aiInput.focus();
  }

  startInterimRecognition() {
    if (!this.speechState.supported || !this.speechState.recognition) return;
    if (this.speechState.isListening) return;

    this.speechState.baseText = this.elements.aiInput.value.trim();
    this.speechState.isListening = true;
    this.speechState.isInterim = true;
    this.speechState.suppressErrors = true;
    this.safeStartRecognition();
  }

  stopInterimRecognition() {
    if (!this.speechState.isInterim || !this.speechState.recognition) return;
    this.speechState.isListening = false;
    this.speechState.isInterim = false;
    this.speechState.suppressErrors = false;
    try {
      this.speechState.recognition.stop();
    } catch (error) {
      // Ignore stop errors for interim mode.
    }
  }

  safeStartRecognition() {
    if (!this.speechState.recognition || this.speechState.isStarting) return;
    this.speechState.isStarting = true;
    try {
      this.speechState.recognition.start();
    } catch (error) {
      if (error && error.name !== 'InvalidStateError') {
        console.error('Speech recognition restart failed:', error);
        this.speechState.isListening = false;
        this.setVoiceListening(false);
      }
    } finally {
      setTimeout(() => {
        this.speechState.isStarting = false;
      }, 100);
    }
  }

  showSpeechError(message) {
    const errorMessage = this.createMessage('error', message, 'System');
    this.elements.chatMessages.appendChild(errorMessage);
    this.scrollToBottom();
  }

  async toggleServerRecording() {
    if (!this.mediaState.supported) return;

    if (this.mediaState.isRecording) {
      this.stopMediaRecording();
      return;
    }

    await this.startMediaRecording();
  }

  async toggleServerStreamingRecording() {
    if (!this.mediaState.supported) return;

    if (this.mediaState.isRecording) {
      this.stopStreamingRecording();
      return;
    }

    await this.startStreamingRecording();
  }

  async startMediaRecording() {
    try {
      this.mediaState.isStreaming = false;
      const constraints = this.mediaState.deviceId
        ? { audio: { deviceId: { exact: this.mediaState.deviceId } } }
        : { audio: true };
      const stream = await navigator.mediaDevices.getUserMedia(constraints);
      this.mediaState.stream = stream;
      this.mediaState.baseText = this.elements.aiInput.value.trim();

      const recorderOptions = this.mediaState.mimeType ? { mimeType: this.mediaState.mimeType } : undefined;
      const recorder = new MediaRecorder(stream, recorderOptions);

      this.mediaState.chunks = [];
      recorder.addEventListener('dataavailable', (event) => {
        if (event.data && event.data.size > 0) {
          this.mediaState.chunks.push(event.data);
        }
      });

      recorder.addEventListener('stop', async () => {
        const blobType = recorder.mimeType || this.mediaState.mimeType || 'audio/webm';
        const audioBlob = new Blob(this.mediaState.chunks, { type: blobType });
        this.mediaState.chunks = [];
        this.stopMediaStream();
        await this.transcribeAudio(audioBlob);
      });

      recorder.start();
      this.startSilenceMonitor(stream);
      this.mediaState.recorder = recorder;
      this.mediaState.isRecording = true;
      this.setVoiceListening(true);
      this.elements.aiInput.focus();
    } catch (error) {
      console.error('Microphone access error:', error);
      this.setVoiceListening(false);
    }
  }

  stopMediaRecording() {
    if (!this.mediaState.recorder) return;
    this.mediaState.isRecording = false;
    this.setVoiceListening(false);
    this.mediaState.recorder.stop();
  }

  async startStreamingRecording() {
    const token = this.getAuthToken();
    if (!token) {
      this.checkAuthStatus();
      return;
    }

    try {
      this.mediaState.isStreaming = true;
      const constraints = this.mediaState.deviceId
        ? { audio: { deviceId: { exact: this.mediaState.deviceId } } }
        : { audio: true };
      const stream = await navigator.mediaDevices.getUserMedia(constraints);
      this.mediaState.stream = stream;
      this.mediaState.baseText = this.elements.aiInput.value.trim();
      this.mediaState.pendingChunks = [];
      this.mediaState.lastStreamText = '';

      const socket = this.createSpeechSocket(token);
      this.mediaState.streamSocket = socket;

      const recorderOptions = this.mediaState.mimeType ? { mimeType: this.mediaState.mimeType } : undefined;
      const recorder = new MediaRecorder(stream, recorderOptions);

      recorder.addEventListener('dataavailable', (event) => {
        if (!event.data || event.data.size === 0) return;
        this.queueSpeechChunk(event.data);
      });

      recorder.addEventListener('stop', () => {
        this.finishStreamingSession();
      });

      recorder.start(400);
      this.startSilenceMonitor(stream);
      this.mediaState.recorder = recorder;
      this.mediaState.isRecording = true;
      this.setVoiceListening(true);
      this.startInterimRecognition();
      this.elements.aiInput.focus();
    } catch (error) {
      console.error('Microphone access error:', error);
      this.setVoiceListening(false);
    }
  }

  stopStreamingRecording() {
    if (!this.mediaState.recorder) return;
    this.mediaState.isRecording = false;
    this.setVoiceListening(false);
    this.stopInterimRecognition();
    this.mediaState.recorder.stop();
  }

  createSpeechSocket(token) {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const wsUrl = `${protocol}://${window.location.host}/ws/speech`;
    const socket = new WebSocket(wsUrl, ['supabase-auth', token]);

    socket.addEventListener('open', () => {
      const startMessage = {
        type: 'start',
        contentType: this.mediaState.mimeType || 'audio/webm',
        filename: 'speech.webm',
        language: 'en-US'
      };
      socket.send(JSON.stringify(startMessage));
      this.flushSpeechQueue();
    });

    socket.addEventListener('message', (event) => {
      this.handleSpeechMessage(event.data);
    });

    socket.addEventListener('close', () => {
      this.mediaState.streamSocket = null;
    });

    socket.addEventListener('error', (error) => {
      console.error('Speech socket error:', error);
    });

    return socket;
  }

  queueSpeechChunk(blob) {
    const socket = this.mediaState.streamSocket;
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(blob);
      return;
    }
    this.mediaState.pendingChunks.push(blob);
  }

  flushSpeechQueue() {
    const socket = this.mediaState.streamSocket;
    if (!socket || socket.readyState !== WebSocket.OPEN) return;
    const queued = this.mediaState.pendingChunks;
    this.mediaState.pendingChunks = [];
    queued.forEach(chunk => socket.send(chunk));
  }

  finishStreamingSession() {
    this.stopMediaStream();
    const socket = this.mediaState.streamSocket;
    if (!socket) return;
    if (socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type: 'stop' }));
    }
  }

  handleSpeechMessage(data) {
    let payload = null;
    try {
      payload = JSON.parse(data);
    } catch (error) {
      return;
    }

    if (payload.type === 'transcript') {
      const base = this.mediaState.baseText.trim();
      const combined = [base, payload.text].filter(Boolean).join(' ').trim();
      this.mediaState.lastStreamText = payload.text;
      this.elements.aiInput.value = combined;
      this.validateInput();
      this.updateCharCounter();
      if (payload.isFinal && combined) {
        this.sendMessage();
      }
      return;
    }

    if (payload.type === 'error') {
      const errorMessage = this.createMessage('error', `Voice input error: ${payload.error}`, 'System');
      this.elements.chatMessages.appendChild(errorMessage);
      this.scrollToBottom();
      this.stopStreamingRecording();
      return;
    }

    if (payload.type === 'done') {
      if (this.mediaState.streamSocket) {
        this.mediaState.streamSocket.close();
        this.mediaState.streamSocket = null;
      }
      return;
    }
  }

  stopMediaStream() {
    this.stopSilenceMonitor();
    this.stopInterimRecognition();
    if (!this.mediaState.stream) return;
    this.mediaState.stream.getTracks().forEach(track => track.stop());
    this.mediaState.stream = null;
  }

  startSilenceMonitor(stream) {
    const AudioContext = window.AudioContext || window.webkitAudioContext;
    if (!AudioContext) return;

    const audioContext = new AudioContext();
    const analyser = audioContext.createAnalyser();
    analyser.fftSize = 2048;
    const source = audioContext.createMediaStreamSource(stream);
    source.connect(analyser);

    this.mediaState.audioContext = audioContext;
    this.mediaState.analyser = analyser;
    this.mediaState.hasSpoken = false;
    this.mediaState.lastSpeechAt = 0;

    const data = new Uint8Array(analyser.fftSize);
    const silenceThreshold = 0.02;
    const silenceTimeoutMs = 1200;

    const monitor = () => {
      if (!this.mediaState.isRecording || !this.mediaState.analyser) return;
      analyser.getByteTimeDomainData(data);
      let sum = 0;
      for (let i = 0; i < data.length; i++) {
        const v = (data[i] - 128) / 128;
        sum += v * v;
      }
      const rms = Math.sqrt(sum / data.length);
      const now = Date.now();

      if (rms > silenceThreshold) {
        this.mediaState.hasSpoken = true;
        this.mediaState.lastSpeechAt = now;
        this.setVoiceSpeaking(true);
      } else if (this.mediaState.hasSpoken) {
        if (now - this.mediaState.lastSpeechAt > silenceTimeoutMs) {
          this.setVoiceSpeaking(false);
          if (this.mediaState.isStreaming) {
            this.stopStreamingRecording();
          } else {
            this.stopMediaRecording();
          }
          return;
        }
        this.setVoiceSpeaking(false);
      }

      this.mediaState.monitorId = requestAnimationFrame(monitor);
    };

    this.mediaState.monitorId = requestAnimationFrame(monitor);
  }

  stopSilenceMonitor() {
    if (this.mediaState.monitorId) {
      cancelAnimationFrame(this.mediaState.monitorId);
      this.mediaState.monitorId = null;
    }
    if (this.mediaState.audioContext) {
      this.mediaState.audioContext.close().catch(() => {});
      this.mediaState.audioContext = null;
    }
    this.mediaState.analyser = null;
    this.mediaState.hasSpoken = false;
    this.mediaState.lastSpeechAt = 0;
    this.setVoiceSpeaking(false);
  }

  async transcribeAudio(blob) {
    const token = this.getAuthToken();
    if (!token) {
      this.checkAuthStatus();
      return;
    }

    try {
      const formData = new FormData();
      formData.append('audio', blob, 'speech.webm');

      const response = await fetch('/api/speech/transcribe', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.details || data.error || 'Failed to transcribe audio');
      }

      const base = this.mediaState.baseText.trim();
      const combined = [base, data.text].filter(Boolean).join(' ').trim();
      this.elements.aiInput.value = combined;
      this.validateInput();
      this.updateCharCounter();
      if (combined) {
        this.sendMessage();
      }
    } catch (error) {
      console.error('Transcription error:', error);
      const errorMessage = this.createMessage('error', `Voice input error: ${error.message}`, 'System');
      this.elements.chatMessages.appendChild(errorMessage);
      this.scrollToBottom();
    }
  }

  setVoiceListening(isListening) {
    if (!this.elements.micButton) return;

    if (isListening) {
      this.elements.micButton.classList.remove('bg-gray-700', 'hover:bg-gray-600');
      this.elements.micButton.classList.add('bg-red-600', 'hover:bg-red-500');
      this.elements.micButton.title = 'Stop voice input';
    } else {
      this.elements.micButton.classList.remove('bg-red-600', 'hover:bg-red-500');
      this.elements.micButton.classList.add('bg-gray-700', 'hover:bg-gray-600');
      this.setVoiceSpeaking(false);
      this.elements.micButton.title = 'Start voice input';
    }
  }

  setVoiceSpeaking(isSpeaking) {
    if (!this.elements.micButton) return;
    if (isSpeaking) {
      this.elements.micButton.classList.add('ring-2', 'ring-red-300', 'shadow-lg', 'shadow-red-500/30');
    } else {
      this.elements.micButton.classList.remove('ring-2', 'ring-red-300', 'shadow-lg', 'shadow-red-500/30');
    }
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

  getUserAvatar() {
    // Try to get avatar from Supabase auth manager
    if (window.authManager && window.authManager.currentUser) {
      const user = window.authManager.currentUser;

      // Try multiple possible avatar fields
      const avatarFields = [
        user.user_metadata?.avatar_url,
        user.user_metadata?.picture,
        user.user_metadata?.avatar,
        user.identities?.[0]?.identity_data?.avatar_url,
        user.identities?.[0]?.identity_data?.picture
      ];

      for (const avatar of avatarFields) {
        if (avatar) {
          return avatar;
        }
      }
    }

    // Default fallback avatar
    return 'https://via.placeholder.com/32/75a7da/FFFFFF?text=U';
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

    if (this.speechState.isListening && this.speechState.recognition) {
      this.speechState.isListening = false;
      this.speechState.recognition.stop();
      this.setVoiceListening(false);
    }

    if (this.mediaState.isRecording) {
      if (this.mediaState.isStreaming) {
        this.stopStreamingRecording();
      } else {
        this.stopMediaRecording();
      }
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
      // Get user avatar from Supabase auth
      const userAvatar = this.getUserAvatar();
      avatar.src = userAvatar;
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
