# Multi-tenancy design

## Goal and recommendation

Use a **pooled schema**: all tenants share the same MySQL tables, and each tenant-owned row is identified by `tenant_id`. Authenticate each request with a tenant API key, resolve that key to a tenant in middleware, place the tenant ID in request context, and require every tenant-owned repository operation to scope its SQL by that ID.

This is the simplest fit for the existing single-database handler → service → repository application. Database-per-tenant can provide stronger isolation, but adds provisioning, connection management, migrations, and operational cost that are not needed for the initial isolation goal.

This design is implemented as a pooled-schema proof of concept. It intentionally avoids rate limiting, metering, database-per-tenant, and a general migration framework.

## Schema

The schema uses `cars`, `attribute_status`, and `health_events`.

The `tenants` table has an internal numeric primary key, tenant name, unique API-key hash, and creation time. Only the hash of a randomly generated high-entropy API key is stored; the plaintext key is shown once at provisioning time.

Each tenant-owned application table has a non-null `tenant_id`:

- `cars`: use `(tenant_id, vin)` as the primary key.
- `attribute_status`: the checked schema stores all metric columns in one row per car, so `(tenant_id, vin)` is the correct primary key; reference `cars(tenant_id, vin)`.
- `health_events`: the checked schema has no auto-increment ID, so use `(tenant_id, vin, attribute_name, start_time)` as the primary key and reference `cars(tenant_id, vin)`.

The composite foreign keys ensure a child row cannot reference a car belonging to another tenant.

### POC schema setup

For this proof of concept, edit `internal/db/migrations/schema.sql` directly and initialize a **fresh, disposable database**. This is not a migration of the old schema and must not be applied over an existing database containing data. If preserving existing data becomes necessary, stop and write a tested backfill migration before changing that database. The schema creates `car_monitor`, which must match `DB_NAME` in `.env`.

## Request authentication and tenant context

Middleware wraps every tenant-owned API route:

1. Read the `X-API-Key` request header.
2. Hash the supplied key with SHA-256 and find the matching row in `tenants.api_key_hash`.
3. Return `401 Unauthorized` for a missing or unrecognized key.
4. Put the resolved tenant ID in `r.Context()` and pass the request to the next handler.

Never accept tenant identity from a client-controlled tenant ID header, query parameter, URL, or JSON field. The only trusted tenant ID is the one resolved from the API key. Keep tenant provisioning separate from the tenant-owned API; do not add an unauthenticated public endpoint that can create arbitrary tenants or keys.

The `cmd/provision` CLI creates tenants out of band: it generates a cryptographically random 32-byte API key, stores only its SHA-256 hash, and prints the plaintext key once. Run it only from a trusted environment with database access. Tests can use the same key generation/hash helpers without logging or printing the key.

## Propagate and enforce isolation

Request context flows through the stack:

`http.Request.Context()` → handler → service → repository → `database/sql` context-aware query.

Repository methods require context, fail closed when no tenant ID is present, and use tenant-scoped SQL for **every** read, insert, update, and delete:

- Car listing and all car CRUD operations.
- Telemetry's simulation check and status upsert.
- Current and full status reads.
- Health event creation, closure, and listing.

Include `tenant_id` in inserts and use it in update/delete predicates. For multi-table reads, join on both tenant ID and VIN. In particular, full status/event queries must not rely on VIN alone. Keep the composite foreign keys as a database-level backstop rather than a replacement for repository scoping.

For a tenant-scoped route, a VIN belonging only to a different tenant should behave as not found; do not disclose that another tenant owns it.

## Generator integration

The telemetry generator reads `API_KEY` from the environment and sends it as `X-API-Key` on both requests. The plaintext key is not stored in `config.yaml` or logged. Provision a generator credential for one tenant only.

## Verification checklist

Run integration tests against a disposable MySQL container started with `docker compose up -d mysql`. The compose service initializes the checked-in schema and publishes a separate local test port. Set `MYSQL_TEST_DSN` for the Go test process to opt into database integration tests; tests skip cleanly when it is unset. Use two tenants with different keys:

- Missing and invalid API keys return `401`.
- Tenant A creates a VIN; tenant B cannot retrieve, update, delete, or read its status/events.
- Both tenants can register the same VIN and observe only their own car and data.
- Listing cars as tenant B never returns tenant A's rows.
- Telemetry from one tenant cannot update another tenant's car or create/close its events.
- Child rows with a mismatched `(tenant_id, vin)` are rejected by the database foreign key.
- Repository operations fail if called without tenant context.

The POC schema is for a fresh database only; legacy-data migration is explicitly out of scope.

## Implemented scope

- Fresh-install tenant schema with composite keys and foreign keys.
- API-key hashing, middleware, tenant context propagation, tenant-scoped request logging, and repository SQL.
- Out-of-band CLI provisioning and generator API-key header.
- Unit tests for key handling, middleware, generator headers, and missing tenant contexts.
- Opt-in MySQL integration tests for two-tenant isolation.

Do not add rate limiting, usage metering, or per-tenant database provisioning as part of this initial isolation milestone.
