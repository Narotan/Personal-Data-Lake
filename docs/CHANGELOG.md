# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Additional analytics queries and aggregations
- Data export functionality (CSV, JSON)
- Multi-user support (optional)
- Mobile responsive improvements
- Dark/Light theme toggle

---

## [2.0.0] - 2025-12-09

### Added
- 🎨 **Modern Web Frontend** built with React 18 + TypeScript
    - Interactive dashboard with real-time data visualization
    - Responsive design with TailwindCSS
    - Smooth animations with Framer Motion
- 📊 **Dashboard Components**
    - KPI Cards with trend indicators
    - Productivity Chart (WakaTime coding time visualization)
    - Health/Activity Chart (Google Fit steps tracking with daily goals)
    - Language Distribution (Pie chart with percentages)
    - Top Projects (Bar chart with time breakdown)
    - Top Applications (ActivityWatch usage statistics)
    - Schedule Timeline (Google Calendar events)
- 📅 **Advanced Date Range Picker**
    - Quick presets (Today, Yesterday, Last 7/30 Days, etc.)
    - Custom date range selection
    - Current Month and Last Year views
    - All-time statistics
- 🧠 **Smart Data Aggregation**
    - Daily data for periods ≤ 90 days
    - Automatic monthly aggregation for periods > 90 days
    - Optimized performance for large datasets
- 🔄 **REST API v1** with dedicated endpoints
    - `/api/v1/wakatime/stats` - Coding statistics
    - `/api/v1/googlefit/stats` - Health & fitness data
    - `/api/v1/googlecalendar/events` - Calendar events
    - `/api/v1/activitywatch/stats` - Computer activity
- 🎯 **Enhanced User Experience**
    - Loading states and skeletons
    - Error handling with user-friendly messages
    - Automatic data refresh
    - Setup page for OAuth configuration
    - Auth success page with redirect
- ⚡ **Performance Optimizations**
    - Parallel data fetching with React hooks
    - Efficient SQL queries with aggregations
    - Optimized database views
    - Vite-powered fast development builds

### Changed
- 🔧 Scheduler interval changed from 30 to 10 minutes (configurable)
- 📦 Restructured project with `api/v1` package separation
- 🗄️ Improved database schema with additional indexes
- 📝 Enhanced API response models with computed fields
- 🎨 Updated UI/UX for better data presentation

### Improved
- 📖 Comprehensive README with architecture documentation
- 🏗️ Detailed backend and frontend structure descriptions
- 📚 API documentation with request/response examples
- 🚀 Simplified quickstart guide
- 🐳 Docker Compose improvements

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

---

## 🏷️ How to Create a Release

### Подготовка
1. Убедитесь, что все изменения закоммичены
2. Обновите версию в `CHANGELOG.md`
3. Обновите версию в `web/package.json` (если нужно)
4. Закоммитьте изменения:
```bash
git add .
git commit -m "chore: bump version to 2.0.0"
git push origin master
```

### Создание тега
```bash
# Создать аннотированный тег
git tag -a v2.0.0 -m "Release v2.0.0 - Web Frontend"

# Отправить тег в GitHub
git push origin v2.0.0
```

### Создание GitHub Release
1. Перейти на GitHub: https://github.com/Narotan/Personal-Data-Lake/releases
2. Нажать "Draft a new release"
3. Выбрать созданный тег `v2.0.0`
4. Заполнить:
   - **Release title**: `v2.0.0 - Web Frontend & Dashboard`
   - **Description**: скопировать из CHANGELOG.md раздел `[2.0.0]`
5. Опционально: прикрепить собранные бинарники
6. Нажать "Publish release"

### Автоматизация (опционально)
Можно настроить GitHub Actions для автоматического создания релиза при пуше тега:

```yaml
# .github/workflows/release.yml
name: Release
on:
  push:
    tags:
      - 'v*'
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          body_path: docs/CHANGELOG.md
          draft: false
          prerelease: false
```

---

## Upgrade Guide

### From v1.0.0 to v2.0.0

#### Database
Никаких изменений в схеме БД не требуется. Все миграции совместимы.

#### Backend
```bash
# Остановить старую версию
make stop

# Обновить код
git pull origin master

# Пересобрать
make build

# Запустить
make start
```

#### Frontend (новое)
```bash
cd web
npm install
npm run dev
```

#### Environment Variables
Добавьте в `.env` (если еще нет):
```env
# Фронтенд будет доступен на порту 5173 (Vite dev server)
# Backend API остается на порту 8000
API_BASE_URL=http://localhost:8000
```

#### Порты
- **Backend API**: `8000` (без изменений)
- **Frontend Dev**: `5173` (новый, для разработки)
- **Frontend Prod**: собрать статику и раздавать через Nginx/Apache

#### Системные требования
- **Node.js**: >= 18.0.0 (для фронтенда)
- **npm**: >= 9.0.0
- Остальные требования без изменений
