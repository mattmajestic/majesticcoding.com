# Changelog

## [v1.0.0] - YYYY-MM-DD
### Added
- Initial **Gin-based live streaming platform** setup.
- **Gin Router** configured with:
  - `/live` handler rendering base template with partials (`stream`, `chat`, `footer`).
  - `/api/stream/status` endpoint to check AWS IVS RTMP server status.
  - `/api/stats/:provider` endpoints for YouTube, GitHub, Twitch, and LeetCode.
  - Swagger UI integration at `/swagger`.
- **Template architecture**:
  - Base layout with header, footer, CSS/JS includes.
  - Partials for `stream`, `chat`, and `footer`.
- **Static assets** pipeline:
  - Tailwind CSS and HTMX integration.
  - JS to dynamically query API and update live stream UI.
- **WebSocket hub** for real-time chat.
- **AWS IVS integration**:
  - Terraform module for channel + stream key provisioning.
  - Go backend polling IVS for live status.
- **Mermaid architecture diagrams**:
  - OBS → AWS IVS → Go API → Live Page flow.
  - Gin API + Swagger architecture.
  - UI partials + handler → API → JS update flow.
- Postgres/Neon database connection for storing stats/history.

---
