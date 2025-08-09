// stream.js — HLS player that follows /api/stream/status

const video = document.getElementById("video");
const offline = document.getElementById("offline-message");

// Hardcode for now; swap to server-rendered later if you want
const awsURL = "https://362a285634f5.us-east-1.playback.live-video.net/api/video/v1/us-east-1.854196226787.channel.y8EpIc8CZ44a.m3u8";
const streamSrc = awsURL;
video.muted = true;
video.autoplay = true;
video.playsInline = true;

function showOffline() {
  offline.classList.remove("hidden");
  video.classList.add("hidden");
}
function showVideo() {
  video.classList.remove("hidden");
  offline.classList.add("hidden");
}

let hls = null;
let usingNative = false;

async function isLiveNow() {
  const ctrl = new AbortController();
  const t = setTimeout(() => ctrl.abort(), 2500); // short timeout

  try {
    // 1) Fetch the master (IVS returns master first)
    const masterRes = await fetch(
      streamSrc + (streamSrc.includes("?") ? "&" : "?") + "ping=" + Date.now(),
      { cache: "no-store", signal: ctrl.signal }
    );
    clearTimeout(t);
    if (!masterRes.ok) return false;

    const master = await masterRes.text();
    if (!master.trim().startsWith("#EXTM3U")) return false;

    // Find first non-comment line
    const first = master.split("\n").map(s => s.trim()).find(l => l && !l.startsWith("#"));
    if (!first) return false;

    // 2) If it’s a VARIANT (.m3u8), fetch it; if it’s a SEGMENT (.ts), we’re already live
    if (!first.endsWith(".m3u8")) {
      return true; // already a media playlist with segments → live
    }

    const variantURL = new URL(first, streamSrc).toString();
    const variantRes = await fetch(variantURL + (variantURL.includes("?") ? "&" : "?") + "ping=" + Date.now(), { cache: "no-store" });
    if (!variantRes.ok) return false;

    const variant = await variantRes.text();
    // A media playlist should contain segment lines (non-#). If present → live.
    return variant.split("\n").some(l => l && !l.startsWith("#"));
  } catch {
    return false;
  }
}


function startPlayer() {
  if (window.Hls && Hls.isSupported()) {
    if (!hls) {
      hls = new Hls({ debug: false });
      hls.attachMedia(video);
      hls.on(Hls.Events.MEDIA_ATTACHED, () => hls.loadSource(streamSrc));
      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        showVideo();
        video.play().catch(() => {});
      });
      hls.on(Hls.Events.ERROR, (_, d) => {
        if (d?.fatal) showOffline();
      });
    }
  } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
    if (!usingNative) {
      usingNative = true;
      video.src = streamSrc;
      video.addEventListener("loadedmetadata", () => {
        showVideo();
        video.play().catch(() => {});
      }, { once: true });
      video.addEventListener("error", showOffline);
    }
  } else {
    showOffline();
  }
}

function stopPlayer() {
  if (hls) {
    try { hls.destroy(); } catch {}
    hls = null;
  }
  if (usingNative) {
    video.removeAttribute("src");
    try { video.load(); } catch {}
    usingNative = false;
  }
  showOffline();
}

async function tick() {
  const live = await isLiveNow();
  if (live) startPlayer(); else stopPlayer();
}

document.addEventListener("DOMContentLoaded", () => {
  tick();
  setInterval(tick, 15000); // re-check every 15s
});
