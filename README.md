# Go REST Starter

Backend-only REST API template with **Go, Chi, Clean Architecture, Ed25519 JWT validation**, an in-memory example store, PostgreSQL via pgx, and an optional MongoDB adapter. No frontend included; it integrates with an existing Next.js/Auth.js Google sign-in.

## Requirements

- Go **1.25+** (pgx security updates require modern Go)
- Optional: Docker with Compose for databases
- Optional: Node.js 20+ to generate a **development-only** test JWT

## Quick start — no database required

```bash
cp .env.example .env
go mod tidy
make keygen
# load the environment (bash):
set -a; source .env; set +a
make run
```

Dependencies are resolved by `go mod tidy` on your first run. This package does not ship a generated `go.sum`; commit the resulting `go.sum` after dependency resolution.

**A fresh checkout will refuse to start without a configured verification public key.** `make keygen` creates local `./.local/keys/api-public.pem` and `api-private.pem` (both ignored by Git). Go reads only the public key. The private key is for a local test client / server-side Next.js signing integration; keep it secret.

In a second terminal:

```bash
curl -i http://localhost:8080/health
TOKEN="$(node scripts/dev-token.mjs)"
curl -i http://localhost:8080/api/v1/items -H "Authorization: Bearer $TOKEN"
curl -i http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"First item"}'
curl -i http://localhost:8080/api/v1/items -H "Authorization: Bearer $TOKEN"
curl -i http://localhost:8080/api/v1/me -H "Authorization: Bearer $TOKEN"
```

The sample token lasts five minutes; regenerate it if you see HTTP 401. **Never use the dev signer as your production authentication endpoint.** In production, only the trusted Next.js server signs API JWTs.

## Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/health` | Public | `{"status":"ok"}` |
| GET | `/api/v1/items` | `items:read` | Lists items owned by the authenticated user |
| POST | `/api/v1/items` | `items:write` | JSON `{ "name": "Example" }`, returns 201 |
| GET | `/api/v1/me` | JWT | Principal's `sub` and `iss` |

Responses are JSON; errors use `{"error":"code"}`. The in-memory store is concurrency-safe but ephemeral. No CORS is enabled; the preferred architecture is server-to-server Next.js -> Go.

## PostgreSQL

```bash
make postgres-up
# .env: STORE_DRIVER=postgres
set -a; source .env; set +a
make migrate # uses psql inside the Postgres container
make run
```

Database credentials in Compose and `.env.example` are for **local development only**. In production use a secret manager, TLS and different credentials. Migrations are intentionally explicit and not automatically applied at startup.

## MongoDB (optional)

```bash
make mongo-up
# .env: STORE_DRIVER=mongodb
set -a; source .env; set +a
make run
```

For production MongoDB, configure authentication, TLS, backups and a compound `{owner_id: 1, created_at: -1}` index. Mongo and Postgres are alternatives, not automatic dual-write storage.

## Quality

```bash
make test
make vet
make check
```

## Identity integration

**Existing Next.js/Auth.js remains responsible for Google login and web sessions.** Next.js validates its session server-side and issues a **separate**, short-lived Ed25519 API JWT. Go verifies the signature using `JWT_PUBLIC_KEY_FILE`, `JWT_ISSUER` and `JWT_AUDIENCE`; it does **not** accept Auth.js cookies as API credentials. See [Next.js integration](docs/NEXTJS_AUTHJS.md) and [architecture](docs/ARCHITECTURE.md).

## Starting a new repository

1. Copy/template the repository and rename the `module` in `go.mod`.
2. Replace every `example.com/go-rest-starter` import with your module path.
3. Run `go mod tidy`, `make check`, and generate new keys.
4. Replace example `items` with your own domain; add migrations/adapters accordingly.
5. Configure TLS, rate limiting, key rotation, monitored logging and secret management before exposing production endpoints.

## AI agents

`AGENTS.md` is the canonical policy. Cursor: `.cursor/rules/project.mdc`; Claude: `.claude/CLAUDE.md`; Codex: root `AGENTS.md` (recognized natively), `.codex/` contains no duplicate rules.
