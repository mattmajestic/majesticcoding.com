const lines = [
  "Full stack developer with a BS in Psychology + MS in Market Research, blending problem-solving with consumer data.",
  "Given that I began working in market research, lived in Excel, shuffled into data science with R/Python to automate data products, and now build scalable web apps like this with Go & Docker in GCP.",
  "I enjoy building efficient applications/services to help businesses run more efficiently and smoother. Feel free to reach out if I can help!"
];

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

function typeLine(lineIdx, charIdx = 0) {
  if (lineIdx >= lines.length) return;
  const el = document.getElementById(`type-line-${lineIdx + 1}`);
  if (!el) return;

  el.textContent = lines[lineIdx].slice(0, charIdx);
  showCursor(el);

  if (charIdx < lines[lineIdx].length) {
    const speed = 30 + Math.random() * 30; // variable speed
    setTimeout(() => typeLine(lineIdx, charIdx + 1), speed);
  } else {
    hideCursor(el);
    setTimeout(() => typeLine(lineIdx + 1, 0), 400);
  }
}

document.addEventListener("DOMContentLoaded", () => typeLine(0, 0));

