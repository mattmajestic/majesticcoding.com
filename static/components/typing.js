function showCursor(el) {
  if (!el.querySelector('.cursor')) {
    const cursor = document.createElement('span');
    cursor.className = 'cursor';
    cursor.textContent = '█';
    cursor.style.animation = 'blink 1s steps(1) infinite';
    el.appendChild(cursor);
  }
}

function hideCursor(el) {
  const cursor = el.querySelector('.cursor');
  if (cursor) cursor.remove();
}

function typeLines(lines, prefix = "type-line-") {
  function typeLine(lineIdx, charIdx = 0) {
    if (lineIdx >= lines.length) return;
    const el = document.getElementById(`${prefix}${lineIdx + 1}`);
    if (!el) return;

    el.textContent = lines[lineIdx].slice(0, charIdx);
    showCursor(el);

    if (charIdx < lines[lineIdx].length) {
      const speed = 30 + Math.random() * 30;
      setTimeout(() => typeLine(lineIdx, charIdx + 1), speed);
    } else {
      hideCursor(el);
      setTimeout(() => typeLine(lineIdx + 1, 0), 400);
    }
  }
  typeLine(0, 0);
}

document.addEventListener("DOMContentLoaded", () => {
  fetch('/static/typing.json')
    .then(res => res.json())
    .then(data => {
      const section = window.TYPING_SECTION || "about";
      if (data[section]) typeLines(data[section], "type-line-");
    });
});