// Interactive Infrastructure Visualization

document.addEventListener('DOMContentLoaded', function() {
    let animationActive = false;

    // Service definitions with details
    const services = {
        'browser': {
            name: 'Web Browser',
            description: 'Client-side application serving the main website with React-like interactions',
            tech: ['HTML5', 'CSS3', 'JavaScript', 'Tailwind CSS'],
            connections: ['nginx', 'gin-api']
        },
        'mobile': {
            name: 'Mobile Applications',
            description: 'Mobile access to the platform via responsive web design',
            tech: ['PWA', 'Responsive Design'],
            connections: ['nginx']
        },
        'obs': {
            name: 'OBS Studio',
            description: 'Broadcasting software for live streaming to RTMP server',
            tech: ['RTMP Protocol', 'Video Encoding'],
            connections: ['rtmp']
        },
        'nginx': {
            name: 'Nginx/Cloudflare',
            description: 'Reverse proxy and load balancer with SSL termination',
            tech: ['Nginx', 'Cloudflare CDN', 'SSL/TLS'],
            connections: ['gin-api', 'rtmp']
        },
        'gin-api': {
            name: 'Gin API Server',
            description: 'Main Go application server handling HTTP requests and routing',
            tech: ['Go', 'Gin Framework', 'REST API', 'WebSockets'],
            connections: ['postgres', 'websocket', 'clerk', 'github', 'twitch', 'youtube', 'spotify', 'leetcode', 'geocoding', 'ivs', 's3']
        },
        'websocket': {
            name: 'WebSocket Hub',
            description: 'Real-time bidirectional communication for chat and live updates',
            tech: ['WebSocket Protocol', 'Go Goroutines', 'Message Broadcasting'],
            connections: ['postgres']
        },
        'rtmp': {
            name: 'RTMP Server',
            description: 'Real-Time Messaging Protocol server for live video streaming',
            tech: ['RTMP', 'HLS', 'Video Processing'],
            connections: ['ivs', 's3']
        },
        'clerk': {
            name: 'Clerk Authentication',
            description: 'Third-party authentication service with JWT token management',
            tech: ['JWT', 'OAuth', 'Session Management'],
            connections: []
        },
        'cloudrun': {
            name: 'Google Cloud Run',
            description: 'Serverless containerized deployment platform for the Go application',
            tech: ['Google Cloud Run', 'Container Runtime', 'Auto-scaling', 'HTTPS'],
            connections: []
        },
        'neon': {
            name: 'Neon Database',
            description: 'Serverless PostgreSQL database with branching and auto-scaling',
            tech: ['Neon Serverless', 'PostgreSQL', 'Auto-scaling', 'Branching'],
            connections: []
        },
        'static': {
            name: 'Static File Storage',
            description: 'Local storage for assets, images, and static content',
            tech: ['File System', 'Static Assets'],
            connections: []
        },
        'github': {
            name: 'GitHub API',
            description: 'Integration for repository data, commits, and project information',
            tech: ['REST API', 'Git Integration'],
            connections: []
        },
        'twitch': {
            name: 'Twitch API',
            description: 'Live streaming platform integration for chat and stream data',
            tech: ['Twitch API', 'IRC Protocol'],
            connections: []
        },
        'youtube': {
            name: 'YouTube API',
            description: 'Video platform integration for channel statistics and content',
            tech: ['YouTube Data API v3'],
            connections: []
        },
        'spotify': {
            name: 'Spotify API',
            description: 'Music streaming service integration for currently playing tracks',
            tech: ['Spotify Web API', 'OAuth 2.0'],
            connections: []
        },
        'leetcode': {
            name: 'LeetCode Integration',
            description: 'Coding challenge platform integration for problem-solving stats',
            tech: ['Web Scraping', 'API Integration'],
            connections: []
        },
        'geocoding': {
            name: 'Geocoding Service',
            description: 'Location services for geographic data and mapping features',
            tech: ['Geocoding API', 'Geographic Data'],
            connections: []
        },
        'ivs': {
            name: 'AWS Interactive Video Service',
            description: 'Managed live streaming service for low-latency video delivery',
            tech: ['AWS IVS', 'HLS', 'Low Latency Streaming'],
            connections: []
        },
        's3': {
            name: 'AWS S3 Storage',
            description: 'Object storage for video files, backups, and media assets',
            tech: ['AWS S3', 'Object Storage', 'CDN'],
            connections: []
        }
    };

    // Service categories for organized display

    // Initialize the visualization
    init();

    function init() {
        setupEventListeners();
        setupServiceHovers();
    }

    function setupEventListeners() {
        document.getElementById('animate-flow').addEventListener('click', animateServices);
        document.getElementById('reset-view').addEventListener('click', resetView);
    }

    function setupServiceHovers() {
        const serviceNodes = document.querySelectorAll('.service-node');
        
        serviceNodes.forEach(node => {
            node.addEventListener('mouseenter', (e) => showServiceInfo(e.target));
            node.addEventListener('mouseleave', hideServiceInfo);
            node.addEventListener('click', (e) => highlightService(e.target));
        });
    }

    function showServiceInfo(node) {
        const serviceId = node.dataset.service;
        const service = services[serviceId];
        
        if (!service) return;

        const serviceInfo = document.getElementById('service-info');
        const serviceName = document.getElementById('service-name');
        const serviceDescription = document.getElementById('service-description');
        const serviceTech = document.getElementById('service-tech');

        serviceName.textContent = service.name;
        serviceDescription.textContent = service.description;
        
        serviceTech.innerHTML = service.tech.map(tech => 
            `<span class="inline-block bg-[#1C3D63] bg-opacity-20 dark:bg-[#60a5fa] dark:bg-opacity-20 text-[#1C3D63] dark:text-[#60a5fa] text-xs px-2 py-1 rounded mr-1 mb-1">${tech}</span>`
        ).join('');

        // Use fixed positioning to prevent twitching
        serviceInfo.classList.remove('opacity-0', 'invisible');
        serviceInfo.classList.add('opacity-100', 'visible');
    }

    function hideServiceInfo() {
        const serviceInfo = document.getElementById('service-info');
        serviceInfo.classList.remove('opacity-100', 'visible');
        serviceInfo.classList.add('opacity-0', 'invisible');
    }

    function highlightService(node) {
        const allNodes = document.querySelectorAll('.service-node');
        
        // Reset all nodes
        allNodes.forEach(n => {
            n.style.opacity = '0.4';
            n.style.transform = 'scale(1)';
        });

        // Highlight clicked service
        node.style.opacity = '1';
        node.style.transform = 'scale(1.05)';
        
        // Show service info
        showServiceInfo(node);
    }

    function animateServices() {
        if (animationActive) return;
        
        animationActive = true;
        const button = document.getElementById('animate-flow');
        button.disabled = true;
        button.innerHTML = '<i class="fas fa-spinner fa-spin mr-2"></i>Animating...';

        const allServices = Object.keys(services);
        allServices.forEach((serviceId, index) => {
            setTimeout(() => {
                const node = document.querySelector(`[data-service="${serviceId}"]`);
                if (node) {
                    node.classList.add('service-pulse');
                    setTimeout(() => {
                        node.classList.remove('service-pulse');
                    }, 800);
                }
            }, index * 200);
        });

        setTimeout(() => {
            animationActive = false;
            button.disabled = false;
            button.innerHTML = '<i class="fas fa-play mr-2"></i>Animate Services';
        }, allServices.length * 200 + 1000);
    }

    function resetView() {
        // Reset all service nodes
        const allNodes = document.querySelectorAll('.service-node');
        allNodes.forEach(node => {
            node.classList.remove('service-pulse');
            node.style.opacity = '1';
            node.style.transform = 'scale(1)';
        });

        hideServiceInfo();
    }
});