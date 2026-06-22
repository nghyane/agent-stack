# Backend API

Go + Postgres. Huma (OpenAPI) + sqlc (typed SQL) + goose migrations + scs sessions.

## Run

```sh
cp .env.example .env
make db-up    # Postgres on 5433
make dev      # API + hot reload, migrations auto-apply
```

Requires Go, Docker, `sqlc`, `goose`. See `make help` for all commands.

## Env

| Variable | Purpose |
|----------|---------|
| `DATABASE_URL` | Postgres connection string |
| `CORS_ORIGINS` | Comma-separated frontend origins |

See `.env.example` for defaults.

## Deploy

Build the Dockerfile. Run behind a reverse proxy (Caddy, nginx) that handles TLS.
Set domain and `CORS_ORIGINS` via env. Migrations run on startup.
