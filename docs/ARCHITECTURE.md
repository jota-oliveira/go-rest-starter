# Architecture

Flow: `Chi router -> authentication/authorization middleware -> HTTP handler -> application use case -> repository port -> persistence adapter`.

- Domain models and validation: `internal/domain/item`.
- Application services and their repository ports: `internal/application/items`.
- Inbound adapter: `internal/adapters/httpapi`; token-verifier adapter: `internal/adapters/auth/jwt`.
- Outbound adapters: `internal/adapters/persistence/{memory,postgres,mongodb}`.
- Startup/dependency injection: `cmd/api/main.go`.

The template avoids a generic repository factory or unnecessary abstraction. Swap persistence using `STORE_DRIVER` without changing domain/use cases. Memory is per-process and resets on restart. Postgres requires the included migration. Mongo is schema-less and indexes should be provisioned in production.

## Authorization design

`Principal` contains subject, issuer and scopes. `GET /api/v1/items` requires `items:read`, `POST /api/v1/items` requires `items:write`. Repository queries filter by owner ID. For multi-tenancy introduce tenant IDs and enforce them in all queries and use cases. Do not trust a client-provided `owner_id`.

`JWT_PUBLIC_KEY_FILE` is an Ed25519 public key. Next.js signs JWTs server-side; Go only verifies. The template does not mint tokens in production. A dev signer in `scripts/dev-token.mjs` is provided strictly for testing, using the locally generated key.

## Future third-party integration

Do not register one signing key per API consumer. Integrate an authorization server with registered client identities, JWKS key rotation, fixed trusted issuers and scope policies. OAuth client credentials is for machine-to-machine integrations. Keep verifier behind the application port. Do not implement a full OAuth server inside this API.
