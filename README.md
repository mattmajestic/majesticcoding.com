# majesticcoding.com in Golang

Welcome to majesticcoding.com.  This was built in Golang with a `Gin API with Swagger`, `AWS IVS RTMP` for `Live Streams`, `Websockets` for `Live Chatting`, `Neon Postgres` for relational database, `Vanilla JS` + `Tailwind CSS` in `Partial templates` on the UI.

### Run with Docker 🐳
```
docker compose up
```
### Run with Go 🐹
```
go run .
```

## Technical Documentation

The following are diagrams outlining the structure of the main functionality of this Go application.

### Gin API with Swagger Architecture Diagram

```mermaid
graph TD
  A["Browser / Client"] -->|HTTP / WS| B["Gin Router"]

  subgraph "Gin App"
    direction TB
    B --> C["Templates (Partials)"]
    B --> D["WebSocket Hub"]
    B --> E["Swagger Docs"]
    B --> F["3rd-Party API Stats"]
    B --> G["Stream Status"]
  end

  F --> H["YouTube / GitHub / Twitch / LeetCode"]
  G --> I["AWS IVS RTMP"]

  C --> A
  D --> A
  E --> A

```

### UI Architecture Diagram

```mermaid
sequenceDiagram
  autonumber
  participant U as "User"
  participant H as "Gin Handler (/live)"
  participant T as "Templates (base + partials)"
  participant P as "Live Page"
  participant API as "Go API"
  participant SV as "IVS / DB"

  U->>H: GET /live
  H->>T: Render base with stream/chat/footer
  T-->>P: HTML
  P->>API: fetch "/api/stream/status"
  API->>SV: Query
  SV-->>API: Status + URL
  API-->>P: JSON { live, playbackUrl }
  P->>P: Update #stream and chat
```


### Live Stream Architecture Diagram

```mermaid
sequenceDiagram
    autonumber
    participant OBS as OBS Encoder
    participant IVS as AWS IVS (RTMP + Playback)
    participant API as Go Backend API
    participant JS as Live Stream Page (JS)

    Note over OBS,IVS: Setup: OBS has IVS ingest endpoint + stream key
    OBS->>IVS: RTMPS connect (ingest_endpoint + stream_key)
    IVS-->>OBS: Accepts stream begins ingest
    IVS-->>IVS: Publishes HLS/DASH playback (playback_url)

    Note over JS,API: Frontend checks live status (polling)
    JS->>API: GET /api/stream/status Check 
    API->>IVS: GetStream(AWS_IVS playback URL)
    IVS-->>API: isLive = true/false, playback_url
    API-->>JS: { live, playbackUrl }

    alt Stream is LIVE
        JS->>JS: Swap UI to "LIVE", show player
        JS->>IVS: Player loads playback_url (HLS)
    else Stream is OFFLINE
        JS->>JS: Show "Offline" UI, hide/disable player
    end
```
