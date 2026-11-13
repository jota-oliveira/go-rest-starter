# Agent instructions — Go REST Starter

This file is the **single source of truth** for AI coding assistants. Read it before modifying code.

## Mission
Build small, maintainable, secure Go REST APIs with Chi and Clean Architecture. Prioritize correctness, clarity, testability and simple dependency boundaries over abstraction for its own sake.

## Language and quality
- Write idiomatic Go. Follow Effective Go, standard library conventions, and `gofmt`.
- Keep exported identifiers documented when public. Return explicit errors and wrap with `%w` when useful.
- Prefer composition, concrete types at boundaries, and small interfaces defined by their consumers.
- Pass `context.Context` first for I/O operations. Honor cancellations and timeouts.
- Avoid globals and hidden runtime dependencies; compose dependencies in `cmd/api`.
- Prefer explicit SQL and pgx over ORMs. Use parameterized SQL, transactions for invariants, and migrations for schema changes.
- Use structured `slog`, without credentials, tokens, or sensitive user data.
- Always add or update tests for behavioral changes; run `go test ./...` and `go vet ./...`.

## Architecture / dependency direction
- `internal/domain`: pure entities, invariants and domain errors; no HTTP, SQL, OAuth, JWT, Chi or DB imports.
- `internal/application`: use cases and inbound/outbound ports; depends only on domain and standard library.
- `internal/adapters`: HTTP/Chi handlers, authentication verifiers, memory/PostgreSQL/MongoDB repositories; adapt external technologies to application interfaces.
- `internal/platform`: environment and technical configuration.
- `cmd`: composition root, wiring, lifecycle, graceful shutdown; no business rules.
- Depend inward. Never call handlers from use cases or import persistence packages into domain/application.
- Do not create generic CRUD repositories or base service classes without a real need.

## HTTP and security
- Validate and size-limit all input. Return consistent JSON error envelopes and appropriate HTTP codes.
- Authenticate every protected route and authorize scopes/ownership for each operation; never trust frontend session alone.
- Verify JWT signature using a pinned asymmetric algorithm, configured public key, issuer, audience, exp and sub; reject untrusted algorithms and unknown issuers.
- Do not add login/password/OAuth flows inside Go for the Next.js integration. Next.js/Auth.js handles login; Go is a resource server.
- Keep private signing keys on the Next.js server; never commit secrets. `.env.example` contains placeholders only.
- Don't expose stack traces, SQL errors or secrets to API consumers.
- Enforce per-owner data access in repository queries. Avoid implicit CORS enablement.
- Implement least privilege scopes and defend externally exposed APIs with TLS, rate limits and observability at deployment.

## Database
- Default to memory so a fresh checkout can run without Docker.
- PostgreSQL and MongoDB are adapters; selection happens at startup. Never leak database-specific models into domain.
- Keep migration SQL reviewed and versioned. Never run destructive migrations automatically on startup.
- Add indexes and query constraints appropriate to authorization and access patterns.

## Agent execution checklist
1. Understand the use case and the inward port before coding.
2. Make the smallest coherent change with explicit dependencies.
3. Add success/failure/authorization tests.
4. Run formatter, tests, vet; report what did or did not run.
5. Update README and `.env.example` for user-facing configuration changes.
