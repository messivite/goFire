# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Documentation site (HTML) with GitHub Pages deploy at https://messivite.github.io/goFire/
- semantic-release for automated versioning and releases (conventional commits)
- README badges (Go Report Card, pkg.go.dev, License, GitHub stars, etc.)
- Documentation link in README
- Updating section in README with `go get -u`, `go mod tidy`, `go build`

### Changed

- Module path set to `github.com/messivite/goFire`
- Documentation UI redesigned (header, sidebar, cards, callouts)

## [0.1.0] - 2025-03-09

### Added

- Firebase Auth middleware for Bearer token verification (file path or JSON credentials)
- YAML-based API definition (`api.yaml`) with code generation
- CLI commands: `gofire init`, `gofire setup`, `gofire add`, `gofire gen`, `gofire list`, `gofire deploy`
- Interactive setup for port, Firebase credentials, and Upstash Redis
- Chi router integration for HTTP routing
- Vercel serverless deployment support
- Upstash Redis cache package
- Built-in health and root handlers
- MIT license

### Changed

- Initial release
