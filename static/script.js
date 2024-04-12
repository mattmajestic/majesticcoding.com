const body = document.body;

        // Function to toggle between light and dark themes
        function toggleTheme() {
            if (body.classList.contains("dark")) {
                // Switch to light theme
                body.classList.remove("dark");
            } else {
                // Switch to dark theme
                body.classList.add("dark");
            }
        }

        // Add click event listener to the theme toggle button
        const themeToggle = document.getElementById("theme-toggle");
        themeToggle.addEventListener("click", toggleTheme);