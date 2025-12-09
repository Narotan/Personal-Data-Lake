# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [2.1.0] - 2025-12-09

### Added - Security Features 🔒
- **Token Encryption** - AES-256-GCM encryption for all OAuth tokens
    - New `auth/encryption.go` module with AES-GCM implementation
    - Encrypted storage in `tokens.json` with backward compatibility
    - `ENCRYPTION_KEY` environment variable (32 bytes required)
    - `NewFileTokenStorageFromEnv()` helper function
- **CORS Middleware** - Cross-Origin Resource Sharing protection
    - Configurable allowed origins via `ALLOWED_ORIGINS` env var
    - Support for preflight requests (OPTIONS)
    - Credentials support for authenticated requests
- **Rate Limiting** - DDoS and brute-force protection
    - Token bucket algorithm with per-IP limiting
    - Configurable via `RATE_LIMIT_RPS` and `RATE_LIMIT_BURST`
    - Automatic cleanup of old limiters
    - HTTP 429 response for exceeded limits
- **Docker Security** - Unprivileged container execution
    - Created `appuser` (UID 1000) in Docker container
    - Application runs as non-root user
    - Proper file ownership configuration
- **Secret Management** - Environment-based configuration
    - Removed hardcoded API keys from scripts
    - All secrets via environment variables
    - `.env.example` with comprehensive documentation
- **Security Documentation**
    - New `docs/SECURITY.md` with detailed security guide
    - Updated README.md with security section
    - Configuration examples and best practices

### Changed
- Updated all token storage initialization to use encryption
- Modified `scripts/generate_traffic.sh` to use env variables
- Enhanced `Dockerfile` with security best practices
- Updated middleware stack with CORS and rate limiting

### Security
- Fixed: OAuth tokens stored in plaintext (CVE-worthy)
- Fixed: Missing CORS configuration
- Fixed: No rate limiting (DDoS vulnerability)
- Fixed: Docker container running as root
- Fixed: Hardcoded secrets in scripts

### Dependencies
- Added `golang.org/x/time/rate` for rate limiting

## [2.0.0] - 2025-12-09

### Added
- **Modern Web Frontend** built with React 18 + TypeScript
    - Interactive dashboard with real-time data visualization
    - Responsive design with TailwindCSS
    - Smooth animations with Framer Motion
- **Dashboard Components**
    - KPI Cards with trend indicators
    - Productivity Chart (WakaTime coding time visualization)
    - Health/Activity Chart (Google Fit steps tracking with daily goals)
    - Language Distribution (Pie chart with percentages)
    - Top Projects (Bar chart with time breakdown)
    - Top Applications (ActivityWatch usage statistics)
    - Schedule Timeline (Google Calendar events)
- **Advanced Date Range Picker**
    - Quick presets (Today, Yesterday, Last 7/30 Days, etc.)
    - Custom date range selection
    - Current Month and Last Year views
    - All-time statistics
- **Smart Data Aggregation**
    - Daily data for periods ≤ 90 days
    - Automatic monthly aggregation for periods > 90 days
    - Optimized performance for large datasets
- **REST API v1** with dedicated endpoints
    - `/api/v1/wakatime/stats` - Coding statistics
    - `/api/v1/googlefit/stats` - Health & fitness data
    - `/api/v1/googlecalendar/events` - Calendar events
    - `/api/v1/activitywatch/stats` - Computer activity
- **Enhanced User Experience**
    - Loading states and skeletons
    - Error handling with user-friendly messages
    - Automatic data refresh
    - Setup page for OAuth configuration
    - Auth success page with redirect
- **Performance Optimizations**
    - Parallel data fetching with React hooks
    - Efficient SQL queries with aggregations
    - Optimized database views
    - Vite-powered fast development builds

### Changed
- Scheduler interval changed from 30 to 10 minutes (configurable)
- Restructured project with `api/v1` package separation
- Improved database schema with additional indexes
- Enhanced API response models with computed fields
- Updated UI/UX for better data presentation

### Improved
- Comprehensive README with architecture documentation
- Detailed backend and frontend structure descriptions
- API documentation with request/response examples
- Simplified quickstart guide
- Docker Compose improvements

### Technical
- **Frontend Stack**
    - React 18.2 with TypeScript 5.2
    - Vite 5.0 for blazing-fast builds
    - Recharts 2.10 for charts
    - Axios for API communication
    - date-fns for date manipulation
    - lucide-react for icons
- **Backend Enhancements**
    - Structured API versioning (v1)
    - Type-safe handlers with proper error handling
    - Enhanced logging for API requests
    - CORS configuration for frontend

---

## [1.0.0] - 2025-11-04

### Added
- Initial project setup
- REST API with endpoints for all data sources
- **WakaTime integration**
    - OAuth2 authentication
    - Data collection and storage
    - Statistics endpoint
- **Google Fit integration**
    - OAuth2 authentication
    - Steps, calories, distance tracking
    - Statistics endpoint
- **Google Calendar integration**
    - OAuth2 authentication
    - Event collection and storage
    - Events endpoint
- **ActivityWatch integration**
    - Event submission endpoint
    - Statistics endpoint
- **Scheduler** for automatic data collection (every 30 minutes)
- **API Key authentication middleware**
- **PostgreSQL database with migrations**
    - User management
    - WakaTime data schema
    - Google Fit data schema
    - Google Calendar data schema
    - ActivityWatch data schema
- **Structured logging** with zerolog
- **Monitoring stack** (Prometheus + Grafana + Loki)
    - Pre-configured Grafana dashboard
    - Prometheus metrics collection
    - Loki log aggregation
- **Docker support**
    - Dockerfile for application
    - Docker Compose for PostgreSQL
    - Docker Compose for monitoring stack
- **SQLC** for type-safe database queries
- **Environment-based configuration**
- **SystemD services** for ActivityWatch client
- **Documentation**
    - README with setup instructions
    - API documentation
    - Contributing guidelines
    - Project status document
- **Build tools**
    - Makefile for common tasks
    - Setup script for initial configuration
    - Build scripts for ActivityWatch client
- **MIT License**

### Security
- API Key authentication for all endpoints
- OAuth2 tokens stored securely in `tokens.json`
- Sensitive data excluded from git (`.gitignore`)
- Constant-time comparison for API keys

---

## Historical Releases
- **2.0.0** — Major update with web frontend
- **1.0.0** — First public release candidate

