// Terminal Widget for Index Page
(function() {
    function initTerminalWidget() {
        const commands = [
            './about', 
            './docs', 
            './live', 
            './youtube',
            './github',
            './linkedin',
            './twitch',  
            './leetcode'
        ];
        
        let currentIndex = 0;
        const commandElement = document.getElementById('typing-command');
        const cursorElement = document.getElementById('typing-cursor');
        
        if (!commandElement || !cursorElement) return;
        
        function typeCommand(command) {
            commandElement.textContent = '';
            cursorElement.style.display = 'none';
            let charIndex = 0;
            
            function typeNextChar() {
                if (charIndex < command.length) {
                    commandElement.textContent += command[charIndex];
                    charIndex++;
                    setTimeout(typeNextChar, 150); // 150ms per character
                } else {
                    // Show cursor and wait longer before next command
                    cursorElement.style.display = 'inline';
                    setTimeout(() => {
                        currentIndex = (currentIndex + 1) % commands.length;
                        typeCommand(commands[currentIndex]);
                    }, 6000); // Slower - 6 seconds instead of 4
                }
            }
            
            typeNextChar();
        }
        
        // Start animation after page loads
        setTimeout(() => typeCommand(commands[currentIndex]), 1500);
    }

    // Initialize when DOM is ready
    document.addEventListener("DOMContentLoaded", initTerminalWidget);
})();