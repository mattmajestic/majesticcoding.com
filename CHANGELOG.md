# Changelog

## [v2.0.0] - 2025-09-10
### Added
- **Globe widget** with real-time geocoding checkins
  - Interactive 3D globe displaying recent location checkins
  - PostgreSQL persistence with automatic 8-hour cleanup
  - Google Geocoding API integration for city/country data
- **Spotify integration** with OAuth2 flow
  - Real-time "now playing" widget with album artwork
  - Database token persistence with automatic refresh
  - Connected terminal-style UI with progress bar
- **Enhanced UI components**
  - Realistic lava lamp animation with glowing effects
  - Terminal-style widgets with consistent theming
  - Responsive layout with connected widget designs
- **Robust API endpoints**
  - `/api/geocode` - geocoding with automatic checkin storage
  - `/api/checkins/recent` - recent location data (8 hours)
  - `/api/spotify/*` - complete OAuth flow and current track
- **Database schema** with migrations
  - Checkins table with location and geocoded data
  - Spotify tokens table with refresh capability
  - Automatic column additions via ALTER TABLE

### Enhanced
- **Authentication flow** moved from file-based to database storage
- **Error handling** with comprehensive logging and fallbacks
- **Real-time updates** with smart polling and change detection
- **Widget connectivity** with matching heights and visual cohesion

### Technical
- **PostgreSQL** with DOUBLE PRECISION for coordinates
- **Manual OAuth2** implementation for better control
- **WebSocket** real-time features for chat and updates
- **Template partials** for modular UI components

## [v1.0.0] - Initial Release
### Core Platform
- **Gin-based live streaming** with AWS IVS integration
- **WebSocket chat** with real-time messaging
- **API integrations** for GitHub, YouTube, Twitch, LeetCode stats
- **Template architecture** with base layouts and partials
- **Static assets** with Tailwind CSS and HTMX
- **Swagger documentation** for all endpoints
- **PostgreSQL/Neon** database foundation

---
