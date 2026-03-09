# GoFire

[![Go Report Card](https://goreportcard.com/badge/github.com/messivite/goFire)](https://goreportcard.com/report/github.com/messivite/goFire)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/messivite/goFire.svg)](https://pkg.go.dev/github.com/messivite/goFire)
[![Go version](https://img.shields.io/github/go-mod/go-version/messivite/goFire)](https://github.com/messivite/goFire)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/messivite/goFire)](https://github.com/messivite/goFire/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/messivite/goFire)](https://github.com/messivite/goFire/network)
[![GitHub issues](https://img.shields.io/github/issues/messivite/goFire)](https://github.com/messivite/goFire/issues)
[![GitHub last commit](https://img.shields.io/github/last-commit/messivite/goFire)](https://github.com/messivite/goFire/commits/main)
[![GitHub release](https://img.shields.io/github/v/release/messivite/goFire)](https://github.com/messivite/goFire/releases)
[![GitHub contributors](https://img.shields.io/github/contributors/messivite/goFire)](https://github.com/messivite/goFire/graphs/contributors)

A Go toolkit for building Firebase-authenticated APIs with code generation and Vercel deployment.

## Features

- **Firebase Auth middleware** – Bearer token verification (file path or JSON credentials)
- **api.yaml** – Define endpoints in YAML, generate Go handlers and routes
- **CLI** – `gofire add`, `gofire gen`, `gofire list`, `gofire deploy`
- **Vercel-ready** – Serverless deployment with one command
- **Chi router** – Lightweight, idiomatic Go HTTP routing

## Quick Start

```bash
# Install the CLI
go install github.com/messivite/goFire/cmd/gofire@latest

# Initialize a new project
gofire init

# Configure (interactive: port, Firebase, Redis)
gofire setup

# Add endpoints
gofire add endpoint "GET /users"
gofire add endpoint "POST /users --auth"
gofire add endpoint "GET /users/:id --auth"

# Generate handlers and server
gofire gen

# Run locally
go run ./cmd/server
```

## Installation

```bash
go get github.com/messivite/goFire
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `gofire init` | Create default `api.yaml` with health endpoints |
| `gofire setup` | Interactive config (port, Firebase, Redis) and save to `.env` |
| `gofire add endpoint "METHOD /path" [--auth]` | Add an endpoint to `api.yaml` |
| `gofire gen` | Generate handler stubs and server routes from `api.yaml` |
| `gofire list` | List all endpoints from `api.yaml` |
| `gofire deploy` | Interactive Vercel deploy (preview or production) |

## api.yaml

All endpoints are defined in `api.yaml`. Edit it directly or use the CLI.

```yaml
version: "1"
basePath: /api

endpoints:
  - method: GET
    path: /api
    handler: Health
    auth: false
  - method: GET
    path: /api/health
    handler: Health
    auth: false
  - method: GET
    path: /users
    handler: ListUsers
    auth: true
  - method: POST
    path: /users
    handler: CreateUsers
    auth: true
```

- `auth: true` – route is protected by Firebase Auth middleware
- `auth: false` – route is public
- `handler` – generated function name in `handlers/` directory

After editing `api.yaml`, run `gofire gen` to regenerate code.

## Configuration

| Variable | Required | Description |
|----------|----------|-------------|
| `PORT` | No | Server port (default: 8080) |
| `FIREBASE_CREDENTIALS_PATH` | No | Path to Firebase service account JSON (local) |
| `FIREBASE_CREDENTIALS_JSON` | No | Full Firebase credentials JSON string (Vercel) |
| `UPSTASH_REDIS_REST_URL` | No | Upstash Redis REST URL (for cache) |
| `UPSTASH_REDIS_REST_TOKEN` | No | Upstash Redis REST token |

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

## Setup Questions

`gofire setup` asks the following questions interactively:

| Question | Default | Description |
|----------|---------|-------------|
| Server port | `8080` | Local server port |
| Firebase credentials JSON path | (empty) | Path to service account JSON, e.g. `./service-account.json`. Leave empty to skip auth |
| Enable Redis cache (Upstash)? | `n` | `y` or `n` |
| Upstash Redis REST URL | — | Only if Redis enabled, e.g. `https://your-db.upstash.io` |
| Upstash Redis REST Token | — | Only if Redis enabled |
| Save configuration to .env file? | `n` | Writes answers to `.env` |

## Firebase Auth

Optional. To enable:

1. Download a service account JSON from Firebase Console
2. Set `FIREBASE_CREDENTIALS_PATH=./service-account.json`

When enabled, endpoints with `auth: true` require an `Authorization: Bearer <token>` header.

On Vercel, use `FIREBASE_CREDENTIALS_JSON` with the full JSON string.

## Upstash Redis (optional)

Use the `cache` package when `cfg.RedisEnabled` is true:

```go
import "github.com/messivite/goFire/cache"

if cfg.RedisEnabled {
    c, _ := cache.NewUpstashCache(cfg.UpstashRedisRestURL, cfg.UpstashRedisRestToken, "myapp:")
    // Inject into handlers or services
}

// In handler: c.Get(ctx, key), c.SetAsync(ctx, key, data, ttlSeconds)
```

## Vercel Deployment

```bash
gofire deploy
```

Or manually:

```bash
vercel --prod
```

Add environment variables:

```bash
vercel env add FIREBASE_CREDENTIALS_JSON production
```

## Project Structure

```
your-project/
├── api.yaml                 # Endpoint definitions (source of truth)
├── api/
│   └── index.go             # Vercel serverless entry point
├── cmd/
│   ├── gofire/main.go       # CLI tool
│   └── server/main.go       # Local server entry point
├── config/
│   ├── env.go               # Environment config
│   └── version.go           # Version constant
├── handlers/
│   ├── health.go            # Built-in health handler
│   ├── root.go              # Built-in root page
│   └── *.go                 # Generated handler stubs
├── cache/
│   ├── cache.go             # Cache interface
│   └── upstash.go           # Upstash Redis implementation
├── middleware/
│   └── firebase.go          # Firebase Auth middleware
├── server/
│   └── server.go            # Chi router (generated by gofire gen)
├── internal/
│   ├── scaffold/             # Code generation engine
│   └── yaml/                 # api.yaml parser
├── vercel.json
├── .env.example
└── go.mod
```

## How It Works

```
api.yaml  ──→  gofire gen  ──→  handlers/*.go + server/server.go
                                        │
                              go run ./cmd/server
                                   or
                              gofire deploy (Vercel)
```

1. Define endpoints in `api.yaml`
2. Run `gofire gen` to generate handler stubs and server routes
3. Implement your handler logic in `handlers/*.go`
4. Run locally with `go run ./cmd/server` or deploy with `gofire deploy`

## License

MIT
