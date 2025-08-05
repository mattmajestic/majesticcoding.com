const toggleButton = document.getElementById("chat-toggle");
const chatWidget = document.getElementById("chat-widget");

toggleButton.addEventListener("click", () => {
  chatWidget.classList.toggle("collapsed");
});
