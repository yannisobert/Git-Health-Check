# ghhealth

Audit the health of public GitHub repositories. Get a score out of 100 with category breakdowns, improvement suggestions, evolution charts, and repo comparison.

## Stack

- **Backend**: Go (`net/http`, `cobra`, `godotenv`)
- **Frontend**: React + Vite (TypeScript)
- **Distribution**: single Go binary with embedded frontend (`go:embed`)

## Prerequisites

- Go 1.22+
- Node.js 20+
- Make

## Quick start

```bash
# Copy environment variables
cp .env.example .env

# Install frontend dependencies
cd frontend && npm install && cd ..

# Run backend + frontend in parallel
make dev
```

- Frontend: http://localhost:5173 (proxies `/api/*` to Go)
- Backend: http://localhost:8080

## Commands

| Command | Description |
|---|---|
| `make dev` | Start Go server + Vite dev server |
| `make dev-go` | Start Go server only |
| `make dev-front` | Start Vite dev server only |
| `make build` | Build frontend and compile binary to `bin/ghhealth` |
| `make test` | Run Go tests |
| `make tidy` | Tidy Go module dependencies |

## CLI

```bash
go run ./cmd/ghhealth server    # start HTTP server
go run ./cmd/ghhealth version   # print version
```

## Project structure

```
cmd/ghhealth/          CLI entry point
internal/github/       GitHub REST API client
internal/analyzer/     Health check logic
server/                HTTP server + embedded frontend
frontend/              React + Vite UI
```

## Development workflow

```
main          production (protected)
dev           integration branch
feature/*     new features → PR to dev
fix/*         bug fixes → PR to dev
chore/*       maintenance → PR to dev
```

## License

MIT
