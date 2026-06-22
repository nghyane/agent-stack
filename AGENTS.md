# AGENTS.md — agent-stack

> **BOOTSTRAP** — If `template-api/` and `template-app/` both exist but
> `api/` does not, this is a fresh clone. Before anything else, ask the user
> for the project name, Go module path, and shape (full-stack / backend /
> frontend), then bootstrap directly:
>
> 1. Rename directories: `template-api/` → `api/`, `template-app/` → `app/`
> 2. Replace the old Go module path and every import:
>    `github.com/nghiahoang/template-api` → the new module path
>    (all text files; skip `.git/`, `node_modules/`, `.svelte-kit/`)
> 3. Rewrite identity strings **before** the dir-name replace
>    (otherwise `template-app` in package.json becomes `app` and the regex
>    misses):
>    - package.json `name` → project name
>    - apiTitle constant → `"<Project> API"` (search for `apiTitle` in the
>      server package)
> 4. Global-replace dir-name references: `template-app` → `app`,
>    `template-api` → `api` (this also repoints opencode.json paths)
> 5. Write the project's `AGENTS.md` with these sections:
>    - `## Purpose` — interview the user: "What does this project do?"
>    - `## Structure` — map of directories to their AGENTS.md/docs
>    - `## Constraints` — interview: "Any constraints or limitations?"
>    - `## Decisions` — empty, grows as libraries are chosen
>    - `## Patterns` — empty, grows as features emerge
>    - `## Workflow` — starting-a-feature checklist
>    Every section records principles, never versions/ports/commands.
>    Include the INVARIANT note at the top. After writing it, read it
>    back — the AGENTS.md is now the project's operating manual.
> 6. Remove the starter kit's git history and verify the backend builds
>
> If backend/frontend only, delete the dropped side and remove its entry
> from opencode.json `references` after step 4. The frontend's contract
> snapshot is stale — regenerate it from the backend.

This repo is **agent-stack**: templates plus opencode config that make agents
generate architecture-consistent code. It is a source, not a workspace — clone
out, never develop products inside it.

## What's here

- Go backend template: Huma + sqlc + Postgres. Owns the OpenAPI contract.
  Conventions documented in its own AGENTS.md.
- SvelteKit frontend template: Svelte 5 runes, Cloudflare Pages. Official
  Svelte MCP server configured in opencode.json.
- opencode config: project wiring and per-directory AGENTS.md references.
- AGENTS.md template: section-based project AGENTS.md the agent fills in.

## Rules

- Keep templates generic and runnable. No product-specific code.
- Stack changes → update the relevant AGENTS.md (principles, not version pins).
- INVARIANT: AGENTS.md files record principles and where to look — never versions,
  ports, or commands. The agent reads source for volatile details.
