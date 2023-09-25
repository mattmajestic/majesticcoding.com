document.addEventListener("DOMContentLoaded", function () {
    const body = document.body;
    const themeToggle = document.getElementById("theme-toggle");

    // Function to toggle the theme
    function toggleTheme() {
        body.classList.toggle("dark");
        if (body.classList.contains("dark")) {
            themeToggle.innerText = "Light";
            localStorage.setItem("theme", "dark");
        } else {
            themeToggle.innerText = "Dark";
            localStorage.setItem("theme", "light");
        }
    }

    // Check the current theme and set it
    const savedTheme = localStorage.getItem("theme");
    if (savedTheme === "dark") {
        body.classList.add("dark");
        themeToggle.innerText = "Light";
    } else {
        themeToggle.innerText = "Dark";
    }

    // Toggle theme when the button is clicked
    themeToggle.addEventListener("click", function () {
        toggleTheme();
    });
});
