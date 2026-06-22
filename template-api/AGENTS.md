# AGENTS.md — template-api

<!-- INVARIANT: this file records principles and where to look. No version
     numbers, port numbers, or CLI commands — the agent reads source for those. -->

Go backend. Self-hosted (Railway/VPS). A separate frontend consumes this API
and generates its client from the exported OpenAPI spec, so the spec is the
contract between repos.

Stack identity: Go + Chi + Huma (OpenAPI + validation) + sqlc (no ORM) + pgx +
Postgres + goose migrations + scs cookie sessions. Read `go.mod`, `Makefile`,
and `sqlc.yaml` for exact versions, commands, and config — do not assume them.

## Architecture: vertical slices

Code is organised by **feature**, not by layer. Each feature is a self-contained
folder under `internal/feature/`; shared infrastructure lives under
`internal/platform/`. To add capability you add a feature folder — you do not
edit files spread across handler/service/db directories.

```
internal/
  platform/            cross-cutting infra (NOT a feature)
    config/  database/  session/  server/
  feature/
    user/              reference implementation of the slice pattern
      queries.sql      sqlc source
      data/            sqlc-generated (DO NOT EDIT)
      response.go      public response types (feature-qualified) + mapping
      service.go       business logic; transport-agnostic (no huma/http)
      handler.go       Huma I/O + validation + RegisterRoutes + error mapping
      *_test.go
    health/            the smallest slice: just a handler
db/migrations/         shared schema (goose, sequential, global)
```

A feature with no business logic (like `health`) can skip `service.go`. Add a
service only where real logic exists; never a pass-through.

## How to find things

- Commands: `make help` lists every task. Use it instead of guessing.
- Features: each folder under the feature directory owns its full stack.
- Where features are wired: `mountFeatures` in the server package.
- Generated code: each feature's generated data package (from sqlc).
- Migrations: the shared migrations directory (goose, sequential, global).

## Rules that don't change

- **Organise by feature, not layer.** New capability = new folder under
  `internal/feature/`, mounted once in `mountFeatures`.
- **`user` is a pattern to learn, not a folder to clone.** Reproduce the slice
  shape for the new domain; never copy `user/` and rename. Cloning carries over
  logic the new feature doesn't have.
- **Response types are feature-qualified** (`OrderResponse`, not `DTO`). The
  OpenAPI schema namespace is flat, so generic names collide across features.
- **Each feature exposes `NewHandler(...)` + `RegisterRoutes(api)`.** The server
  treats every feature through the same interface.
- **Three layers inside a slice**: `handler` (transport) → `service` (logic) →
  `data` (sqlc). Handlers stay thin; logic lives in the service.
- **Services are transport-agnostic.** They return sentinel errors
  (`var ErrX = errors.New(...)`), never `huma.ErrorXXX`. The handler maps those
  to HTTP in one place (`toHTTP`). Services never import huma or net/http.
- **Never edit generated code** (`data/`). Change the SQL and regenerate.
- **Validation lives in Huma struct tags**, not hand-written `if` checks. The
  same tags drive runtime validation and the OpenAPI contract.
- **DTOs never expose secrets** (e.g. password hashes). Map DB rows to an
  explicit response struct, never return DB models directly.
- **The OpenAPI spec is the cross-repo contract.** Any change to a request or
  response shape must be reflected in the spec the frontend regenerates from.

## Workflow: adding a feature (or an endpoint in one)

Learn the pattern from the `user` slice and reproduce it for the new domain —
do not copy the folder. Cloning drags in user-specific logic (password hashing,
sessions) that the new feature doesn't need. Write each file fresh, following
the same shape:

1. Create `internal/feature/<name>/`. For an endpoint in an existing feature,
   work within that folder.
2. Write SQL in the feature's `queries.sql` (add a migration in
   `db/migrations/` if the schema changes).
3. Regenerate the typed `data/` package (see `make help`).
4. Add `service.go` logic with sentinel errors (skip if there is none).
5. Add `response.go` with feature-qualified response types (e.g.
   `OrderResponse`, not a generic `DTO`) plus the row→response mapping.
6. Add the Huma I/O structs + handler; put validation in struct tags; map
   errors to HTTP in one place (the feature's `toHTTP`).
7. Mount the feature in `mountFeatures` (new feature only).
8. Re-export the OpenAPI spec so the frontend client stays in sync.
9. Verify: build, then test.

## Auth model

Cookie sessions (scs), backed by Postgres. The session platform package
wraps every request; features read the current user through the session
manager — they do not parse cookies or tokens directly. Cookies are HttpOnly;
production tightens them for cross-site use by the frontend.

## Configuration: two tiers

- **Static config** (from env, in the config platform package): wiring + secrets,
  read once at startup, immutable for the process. Connection strings, API
  keys, CORS origins. Changing it means a redeploy/restart — that is correct.
- **Dynamic settings** (DB-backed, in the settings platform package): operational
  toggles that change WITHOUT a restart. Stored as a single JSONB row; reads
  are a lock-free in-memory snapshot; writes persist and fan out to every
  instance via Postgres LISTEN/NOTIFY, with a reload-on-(re)connect to catch
  missed notifications.

Rules:
- **Never put secrets or wiring in dynamic settings.** Secrets stay in env.
- **Read settings via the cached snapshot** — never query the row per request.
- **Add a setting = a new field + a default.** No migration: it is JSONB.
- **LISTEN/NOTIFY needs a direct Postgres connection.** A transaction-mode
  pooler (e.g. pgbouncer in transaction pooling) silently breaks it. Use a
  session-mode/direct connection, or fall back to TTL polling there.

## Deploy & security

- Container build and compose files are in the repo root; the reverse proxy
  handles TLS. Set the domain and CORS origins via env/config, never hard-coded.
- The API is internet-exposed. **CORS is a strict allowlist** (credentials mode
  forbids wildcards). Always set real frontend origins before deploying.
- Migrations run automatically on startup.
