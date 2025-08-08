// Live Stream Stream Player

const video = document.getElementById("video");
const offline = document.getElementById("offline-message");

const streamName = window.STREAMING_KEY || "test123";
const streamSrc = `http://localhost:8081/${streamName}.m3u8`;

// Hidden if Offline
function showOffline() {
  offline.classList.remove("hidden");
  video.classList.add("hidden");
}

function showVideo() {
  video.classList.remove("hidden");
  offline.classList.add("hidden");
}

if (Hls.isSupported()) {
  const hls = new Hls();
  hls.loadSource(streamSrc);
  hls.attachMedia(video);

  hls.on(Hls.Events.ERROR, (event, data) => {
    if (data.fatal) showOffline();
  });

  hls.on(Hls.Events.MANIFEST_PARSED, (event, data) => {
    hls.currentLevel = data.levels.length - 1;
    showVideo();
});

} else if (video.canPlayType("application/vnd.apple.mpegurl")) {
  video.src = streamSrc;
  video.addEventListener("loadedmetadata", showVideo);
  video.addEventListener("error", showOffline);
} else {
  showOffline();
}


