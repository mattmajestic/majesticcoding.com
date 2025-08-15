(function () {
  function init(id) {
    var root = document.getElementById(id);
    if (!root) return;

    var speed = 4; // px/sec, even slower
    var content = root.querySelector('.marquee__content');
    if (!content) return;

    var items = Array.from(content.children);
    var visibleCount = 4;

    // Set container width to show only 4 items
    var itemWidth = items[0].offsetWidth;
    root.style.width = (itemWidth * visibleCount) + "px";
    content.style.width = (itemWidth * items.length * 2) + "px";

    // Duplicate items for seamless loop
    items.forEach(function (item) { content.appendChild(item.cloneNode(true)); });

    var x = 0, last = performance.now(), raf, totalWidth = itemWidth * items.length;

    function frame(now) {
      var dt = (now - last) / 1000;
      last = now;
      x -= speed * dt;
      if (totalWidth > 0 && -x >= totalWidth) x += totalWidth; // wrap
      content.style.transform = 'translateX(' + x + 'px)';
      raf = requestAnimationFrame(frame);
    }

    function start() {
      raf = requestAnimationFrame(frame);
    }

    // Wait for images to load before starting
    var imgs = Array.prototype.slice.call(content.querySelectorAll('img'));
    var loaded = 0;
    if (imgs.length === 0) { start(); return; }
    imgs.forEach(function (img) {
      if (img.complete) { if (++loaded === imgs.length) start(); }
      else img.addEventListener('load', function () { if (++loaded === imgs.length) start(); }, { once: true });
    });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function () { init('cloud-marquee'); });
  } else {
    init('cloud-marquee');