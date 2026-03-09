# GoFire — *Firebase-authenticated Go APIs from YAML*

> 📖 **[Documentation](https://messivite.github.io/goFire/)** — guides, API reference, deployment


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
[![Documentation](https://img.shields.io/badge/docs-messivite.github.io/goFire-blue)](https://messivite.github.io/goFire/)

A Go toolkit for building Firebase-authenticated APIs with code generation and Vercel deployment.

## Features

- **Firebase Auth middleware** – Bearer token verification (file path or JSON credentials)
- **api.yaml** – Define endpoints in YAML, generate Go handlers and routes
- **CLI** – `gofire add`, `gofire gen`, `gofire list`, `gofire deploy`
- **Vercel-ready** – Serverless deployment with one command
- **Chi router** – Lightweight, idiomatic Go HTTP routing

## Quick Start

**Option A — Create from scratch** (single command):

```bash
go install github.com/messivite/goFire/cmd/gofire@latest
gofire new my-api
cd my-api
go mod tidy
make run
```

**Option B — Existing project** (manual setup):

```bash
mkdir my-api && cd my-api
go mod init my-api
go get github.com/messivite/goFire
gofire init
gofire setup
gofire add endpoint "GET /users"
gofire add endpoint "POST /users" --auth
gofire gen
go mod tidy
go run ./cmd/server
```

> **Custom layout (pkg/server, pkg/handler)?** By default `gofire gen` writes to `server/` and `handlers/`. Define paths once in `.gofire.yaml` or `api.yaml`:
> ```yaml
> # .gofire.yaml (project root) or api.yaml output section
> output:
>   serverDir: pkg/server
>   handlersDir: pkg/handler
> ```
> Then run `gofire gen` without flags. CLI flags `--server-dir` / `--handlers-dir` override config.

> **Build errors?** Run `go mod tidy` to fetch transitive dependencies. Run `go run ./cmd/server` from the project root (not `go run .`).

> **Tip:** Or clone this repo as a template: `git clone https://github.com/messivite/goFire.git my-api && cd my-api`

## Install Demo (Video)

[![GoFire Install & Demo](https://img.youtube.com/vi/c6Xz76uUH38/hqdefault.jpg)](https://youtu.be/c6Xz76uUH38)

## Installation

```bash
go get github.com/messivite/goFire
```

### CLI: command not found?

If `gofire` is not found after install, add Go's bin directory to your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Or run with the full path: `$(go env GOPATH)/bin/gofire init`

## Updating

To update to the latest version:

```bash
go get -u github.com/messivite/goFire
go mod tidy
go build ./...
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `gofire new &lt;name&gt;` | Create new project from scratch (directory, go.mod, api.yaml, handlers, server) |
| `gofire init` | Create `api.yaml` and `cmd/server/main.go` in existing project |

### new vs init — when to use which?

| Scenario | Command |
|----------|---------|
| **Starting from scratch** — you have no project yet | `gofire new my-api` |
| **Existing project** — you already have `go.mod` or an existing Go project | `gofire init` |

- **`gofire new`** creates the project directory, runs `go mod init`, adds the goFire dependency, and generates api.yaml, handlers, and server in one go. Use this when you want a ready-to-run API from nothing.
- **`gofire init`** adds api.yaml and cmd/server/main.go to the current directory. Requires an existing `go.mod` (or you'll get a warning). Run `gofire gen` afterward to generate handlers and server.
| `gofire setup` | Interactive config (port, Firebase, Redis) and save to `.env` |
| `gofire add endpoint "METHOD /path" [--auth]` | Add an endpoint to `api.yaml` |
| `gofire gen [--server-dir DIR] [--handlers-dir DIR]` | Generate handler stubs and server routes. Use flags or `api.yaml` output section for custom paths (e.g. `pkg/server`, `pkg/handler`). |
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
- `output` (optional) – custom paths for `gofire gen`. Can be in `api.yaml` or in `.gofire.yaml` at project root. Resolution: CLI flag > `.gofire.yaml` > `api.yaml` output > default:

```yaml
output:
  serverDir: pkg/server
  handlersDir: pkg/handler
```

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

---

Made with ❤️ by [Mustafa Aksoy](https://www.linkedin.com/in/mustafa-aksoy-87532a385/)
