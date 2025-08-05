let idCounter = 0;

function addNode(type) {
    const canvas = document.getElementById("diagram-canvas");
    const node = document.createElement("div");
    node.className = "absolute bg-gray-700 border border-white p-2 rounded cursor-move";
    node.style.top = `${50 + idCounter * 60}px`;
    node.style.left = "50px";
    node.textContent = `${type.toUpperCase()} Node ${idCounter + 1}`;
    node.draggable = true;
    node.dataset.id = idCounter;
    canvas.appendChild(node);
    idCounter++;
}
