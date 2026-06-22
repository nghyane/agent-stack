# agent-stack

[![License: MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)](.)
[![Svelte](https://img.shields.io/badge/svelte-5-FF3E00?logo=svelte)](.)
[![agent-native](https://img.shields.io/badge/agent--native-✓-6e40c9)](.)

One prompt. One minute. A complete Go + SvelteKit project your AI agent
understands from line zero.

## Use case

You start a new full-stack project. You want a Go API and a SvelteKit frontend
that talk through a typed OpenAPI contract. You want your AI coding agent to
navigate it correctly from day one — following conventions, respecting
architecture, never guessing.

You don't want boilerplate generators that dump code you can't explain. You
don't want to manually rename 30 files and rewrite every import. You don't
want scripts you can't trust.

agent-stack solves this. It ships patterns the agent reads and executes —
not opaque scripts, not stale templates. Clone, bootstrap, interview. Done.

## How it works

```
One prompt  →  agent clones  →  renames everything  →  interviews you  →  writes AGENTS.md  →  ready
```

**1. Clone.** `git clone --depth 1` the repo into your project directory.

**2. Bootstrap.** The agent follows the BOOTSTRAP checklist in `AGENTS.md`:
renames directories, rewrites the Go module path and every import, sets the
app name and API title, repoints config — no scripts, no black boxes. The
agent owns every step.

**3. Interview.** The agent writes the project's `AGENTS.md` from scratch:
purpose, constraints, workflow. As the project grows, the agent adds to
Decisions and Patterns. The AGENTS.md becomes a living knowledge base.

**4. Ready.** `go build` passes. `pnpm run check` has zero errors. Your
agent knows the architecture and conventions. Start building.

## Why this over alternatives

**No scripts.** The agent reads `AGENTS.md` and acts. Every rename, every
replace is visible and verifiable. If something breaks, the agent knows
exactly where.

**No stale docs.** Documents record principles and where to look — never
version numbers, ports, or commands. The agent reads source for volatile
details.

**No copy-paste.** Features are vertical slices (handler → service → data).
The agent learns the pattern from an example and reproduces it for new
domains — never clones a folder full of irrelevant logic.

**No setup drift.** Per-directory AGENTS.md files auto-load when the agent
works in each directory. `api/AGENTS.md` for backend conventions, `app/AGENTS.md`
for frontend. The Svelte MCP server provides live docs and static analysis.

## Stack

- **Go backend** — Huma (OpenAPI 3.1), sqlc (typed SQL), pgx + Postgres,
  goose migrations, scs sessions. Two-tier config: static env for secrets,
  DB-backed dynamic settings with Postgres LISTEN/NOTIFY.
- **SvelteKit frontend** — Svelte 5 runes, Cloudflare Pages. Typed API client
  from OpenAPI spec. Runtime backend URL via `$env/dynamic/public`.
- **Agent tooling** — Svelte MCP server in `opencode.json`, per-directory
  AGENTS.md files the agent auto-loads.

## Quick start

```
Start a new project from https://github.com/nghyane/agent-stack:

1. git clone --depth 1 the repo into a fresh directory
2. Ask me for the project name and Go module path
3. Follow the BOOTSTRAP checklist in the repo's AGENTS.md
```

Paste into any AI coding agent.

## License

MIT — see [LICENSE](LICENSE).
