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
# Analyze a repo and display a colored report
ghhealth check owner/repo
ghhealth check https://github.com/owner/repo

# Flags
ghhealth check owner/repo --json            # output raw JSON
ghhealth check owner/repo --no-color        # disable colors
ghhealth check owner/repo --period monthly  # history period (weekly|monthly)

# Start the HTTP server (serves API + embedded frontend)
ghhealth server

ghhealth version
```

## API

All endpoints return JSON. Add a `GITHUB_TOKEN` in `.env` to raise the rate limit from 60 to 5000 req/h.

| Endpoint | Description |
|---|---|
| `GET /api/analyze?repo=owner/repo` | Full health report (score + checks + suggestions) |
| `GET /api/history?repo=owner/repo&period=weekly` | Score evolution over time |
| `GET /api/compare?repo1=o/r1&repo2=o/r2` | Side-by-side comparison with diff |
| `GET /api/rivals?repo=owner/repo` | Suggested similar repos to compare |
| `GET /api/health` | Server health check |

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
