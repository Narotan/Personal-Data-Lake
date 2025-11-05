# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Web frontend for data visualization
- Additional analytics queries
- Data export (CSV, JSON)
- Multi-user support (optional)

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
- **1.0.0** — First public release candidate

---

## Upgrade Guide
TBD
