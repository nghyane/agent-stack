# Frontend

SvelteKit + Svelte 5 runes. Talks to the Go backend through a typed OpenAPI client.

## Run

```sh
pnpm install    # generates typed client from openapi.yaml
cp .env.example .env
pnpm dev
```

Requires Node, pnpm. Backend should be running first.

## Env

| Variable | Purpose |
|----------|---------|
| `PUBLIC_API_URL` | Backend API base URL, exposed to client code |

## Scripts

`dev` `build` `preview` `check` `gen:api` — see `package.json`.

## Deploy

`pnpm build` → Cloudflare Pages. Set `PUBLIC_API_URL` per environment.
Backend's `CORS_ORIGINS` must include the deployed origin.
