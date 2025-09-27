

# <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg" width="40" height="40" style="vertical-align: middle;"> Majestic Coding

**Full Stack Go Web Application with Live Streaming & AI**

[![Go](https://img.shields.io/badge/Go-1.23.6+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-Web%20Framework-00ADD8?style=flat-square&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Neon-316192?style=flat-square&logo=postgresql&logoColor=white)](https://neon.tech/)
[![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docker.com/)
[![AWS](https://img.shields.io/badge/AWS-IVS%20Streaming-FF9900?style=flat-square&logo=amazon-aws&logoColor=white)](https://aws.amazon.com/ivs/)
[![AI](https://img.shields.io/badge/AI-Multi%20Provider-FF6B6B?style=flat-square&logo=openai&logoColor=white)](#ai-integration)

---

## 📋 Contents

**[🚀 Quick Start](#-quick-start)** • **[✨ Features](#-features)** • **[🏗️ Architecture](#️-architecture)** • **[📁 Structure](#-directory-structure)** • **[🛠️ Tech Stack](#️-tech-stack)** • **[🌐 API](#-api-endpoints)** • **[🔧 Setup](#-environment-setup)** • **[🐳 Deploy](#-deployment)**

## 🚀 Quick Start

<table>
<tr>
<td width="50%">

**🐳 Docker (Recommended)**
```bash
docker compose up
```

</td>
<td width="50%">

**🐹 Go Native**
```bash
go run .
```

</td>
</tr>
</table>

**🔧 Development Setup**
```bash
git clone https://github.com/mattmajestic/majesticcoding.com.git
cd majesticcoding.com && go mod tidy && go run .
```

➡️ **Open:** `http://localhost:8080`

## ✨ Features

<table>
<tr>
<td align="center" width="33%">

**🎥 Live Streaming**
AWS IVS • RTMP • HLS
WebSocket Chat • Analytics

</td>
<td align="center" width="33%">

**🤖 AI Integration**
Claude • GPT • Gemini • Groq
RAG • Vector Embeddings

</td>
<td align="center" width="33%">

**📊 Social Analytics**
GitHub • YouTube • Twitch
LeetCode • Real-time Stats

</td>
</tr>
<tr>
<td align="center">

**🔐 Authentication**
Supabase • JWT • Session Cache
OAuth • Security

</td>
<td align="center">

**🌍 Geographic**
Check-ins • 3D Globe
Geocoding • Locations

</td>
<td align="center">

**📡 API Services**
REST • GraphQL • WebSocket
Swagger • Bronze Schema

</td>
</tr>
</table>

### 🎵 **Spotify Integration** • 🔧 **Content Moderation** • 🚀 **Real-time Everything**

## 🏗️ Architecture

**MVC Pattern** • **Microservice Ready** • **Event-Driven** • **Cloud-Native**

```
Frontend ↔ Gin API ↔ Services ↔ PostgreSQL + Vector DB
    ↕         ↕         ↕
WebSocket   REST    External APIs
```

## 📁 Directory Structure

<details>
<summary><b>📦 Click to expand full structure</b></summary>

```bash
majesticcoding.com/
├── 📦 api/                 # Backend API Layer
│   ├── handlers/           # HTTP controllers
│   ├── services/           # Business logic + integrations
│   ├── models/             # Data structures
│   └── middleware/         # Auth, CORS, etc.
├── 📦 db/                  # Database Layer
│   ├── *.go               # Queries, connections, schemas
│   └── migrations/         # Schema changes
├── 📦 static/              # Frontend Assets
│   ├── components/         # JS modules
│   ├── styles/            # Tailwind CSS
│   └── img/               # Static assets
├── 📦 templates/           # HTML Templates
├── 📄 main.go             # Entry point
├── 📄 docker-compose.yml  # Container orchestration
└── 📄 k8s-go.yaml         # Kubernetes deployment
```

</details>

**🏢 API Layer:** REST handlers + business services
**🗄️ Database:** PostgreSQL + Vector embeddings + Session cache
**🎨 Frontend:** Vanilla JS + Tailwind + Go templates
**⚙️ Infrastructure:** Docker + Kubernetes ready

## 🛠️ Tech Stack

<table>
<tr>
<td align="center" width="25%"><img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg" width="40"><br><b>Go 1.23+</b></td>
<td align="center" width="25%"><img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/postgresql/postgresql-original.svg" width="40"><br><b>PostgreSQL</b></td>
<td align="center" width="25%"><img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/docker/docker-original.svg" width="40"><br><b>Docker</b></td>
<td align="center" width="25%"><img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/javascript/javascript-original.svg" width="40"><br><b>JavaScript</b></td>
</tr>
</table>

- **Backend:** Go + Gin + PostgreSQL/Neon + Swagger
- **Frontend:** Vanilla JS + Tailwind + WebSockets
- **AI:** Claude + GPT + Gemini + Groq + RAG
- **Cloud:** AWS IVS + Supabase + Docker + K8s

## 🌐 API Endpoints

| **Category** | **Endpoint** | **Description** |
|:---:|:---:|:---:|
| 🔐 **Auth** | `POST /api/user/sync` | Sync user data |
| 🎥 **Stream** | `GET /api/stream/status` | Live stream status |
| 🤖 **AI** | `POST /api/llm/` | Chat with AI |
| 📊 **Stats** | `GET /api/stats/{platform}` | Social media analytics |
| 💬 **Chat** | `GET /ws/chat` | WebSocket connection |

**🔗 Full Documentation:** `http://localhost:8080/docs`

## 🔧 Environment Setup

<details>
<summary><b>🔑 Environment Variables (Click to expand)</b></summary>

```bash
# Database
DATABASE_URL=postgresql://user:password@host:5432/database

# Authentication
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key

# AI Providers
ANTHROPIC_API_KEY=your-key
OPENAI_API_KEY=your-key
GEMINI_API_KEY=your-key

# Social APIs
GITHUB_TOKEN=your-token
YOUTUBE_API_KEY=your-key
TWITCH_CLIENT_ID=your-id

# AWS IVS
AWS_IVS_CHANNEL_ARN=your-arn
```

</details>

## 🐳 Deployment

<table>
<tr>
<td align="center" width="33%">

**🔧 Development**
```bash
go run .
# or
docker compose up
```

</td>
</tr>
</table>

**⚡ Database:** Enable `pgvector` extension in Neon • Migrations run automatically

## 📊 Architecture Diagrams

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

## 🤝 Contributing

- **Fork → Clone → Code → PR**
- **📝 License:** MIT 
- **🙏 Thanks:** Go Team, Gin, Neon, Supabase

---

**⭐ Star this repo** • **🐛 Report issues** • **💡 Suggest features**

[🌐 Website](https://majesticcoding.com)
