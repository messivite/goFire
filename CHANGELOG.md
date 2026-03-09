# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.14] - 2026-03-09

## [0.1.13] - 2026-03-09

## [0.1.12] - 2026-03-09

### Added

- Custom output paths for `gofire gen`: `--server-dir` and `--handlers-dir` CLI flags (e.g. `pkg/server`, `pkg/handler`)
- Optional `output` in `.gofire.yaml` (project root) or `api.yaml` — define paths once, `gofire gen` uses them without flags
- README and documentation updated with custom layout guidance for existing projects

### Changed

- `gofire gen` derives handler package name from output directory (e.g. `pkg/handler` → `package handler`)
- `GenerateServer` now accepts `handlersDir` and generates correct import paths for any layout

## [0.1.11] - 2026-03-09

## [0.1.10] - 2026-03-09

## [0.1.9] - 2026-03-09

## [0.1.8] - 2026-03-09

## [0.1.7] - 2026-03-09

## [0.1.6] - 2026-03-09

## [0.1.5] - 2026-03-09

## [0.1.4] - 2026-03-09

## [0.1.3] - 2026-03-09

## [0.1.2] - 2026-03-09

### Added

- `gofire new <name>` creates a complete project from scratch (directory, go mod init, api.yaml, handlers, server, .gitignore, Makefile)
- `gofire init` now creates `cmd/server/main.go` in addition to `api.yaml` (module path from `go.mod`)
- `gofire gen` generates `handlers/health.go` and `handlers/root.go` in the project (no manual main entry point needed)
- Documentation site (HTML) with GitHub Pages deploy at https://messivite.github.io/goFire/
- GitHub Pages workflow for docs deployment
- README badges (Go Report Card, pkg.go.dev, License, GitHub stars, etc.)
- Documentation link in README
- Updating section in README with `go get -u`, `go mod tidy`, `go build`
- Install Demo video linked from README and docs home
- Social header icons and footer linking to GitHub, npm, and LinkedIn
- "new vs init" guidance in README and docs

### Changed

- `gofire gen` server template imports handlers from user's module path (e.g. `my-api/handlers`) instead of `github.com/messivite/goFire/handlers`
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
