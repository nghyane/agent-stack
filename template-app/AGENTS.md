# Frontend

<!-- INVARIANT: principles only. No versions, ports, or commands — the agent
     reads source for those. Official MCP server handles live Svelte docs. -->

SvelteKit + Svelte 5 runes, deployed to Cloudflare Pages.

## Official tooling

The Svelte MCP server (`@sveltejs/mcp`) is configured in `opencode.json`.
It provides `get-documentation` for current docs and `svelte-autofixer` for
static analysis. Use it for all Svelte/SvelteKit work.

## Conventions

- **Typed client only.** Every API call goes through the OpenAPI-generated
  client. Regenerate from the OpenAPI spec when the backend changes. Never
  hand-write `fetch` URLs or edit `schema.d.ts`.
- **Runtime backend URL.** `PUBLIC_API_URL` from `$env/dynamic/public`. One
  build targets dev and prod.
- **Client-side SPA.** Cloudflare Pages adapter with static fallback. No
  server endpoints unless explicitly added via `+page.server.ts`.
