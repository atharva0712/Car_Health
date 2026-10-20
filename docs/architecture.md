# Project architecture

## Overview

Car Health is a small Go HTTP service backed by MySQL. It stores registered cars, their latest telemetry-derived status, and status-transition events. A separate process generates sample telemetry and communicates with the backend over HTTP.

## Components

### Backend

`cmd/backend/main.go` wires a standard-library `http.ServeMux`, repositories, services, and handlers. The backend exposes car-management, telemetry, status, and event endpoints.

The application code is organized in layers:

1. **Handlers** in `internal/handler` decode HTTP requests, call services, and construct responses.
2. **Services** in `internal/service` apply business logic, such as classifying metrics and creating or closing health events.
3. **Repositories** in `internal/repository` issue SQL queries against MySQL.
4. **Models** in `internal/models` represent cars, telemetry payloads, statuses, and health events.

`internal/db/connection.go` reads database settings from environment variables (optionally loaded from `.env`) and opens the MySQL connection pool.

### Telemetry generator

`cmd/generator` loads `config.yaml` and starts `internal/generator.Ticker`. On each interval, it calls `GET /cars`, generates values for cars whose simulation is enabled, and submits them to `POST /telemetry`.

The generator is a separate client process. It reads `API_KEY` from the environment and sends it as `X-API-Key` on both requests so it operates within one tenant.

### Configuration

- `.env.example` documents database connection, backend port, and generator API-key environment variables.
- `config.yaml` controls telemetry generation cadence, normal/warning/critical distribution, metric thresholds, and the generator's backend URL.
- The database schema is currently bootstrapped from `internal/db/migrations/schema.sql`. Despite the directory name, this is a schema script rather than an ordered/versioned migration system.

## Data model

The schema creates a tenant table and three tenant-owned application tables:

- **`tenants`**: tenant name and unique SHA-256 API-key hash.
- **`cars`**: one row per tenant and VIN, including whether simulation is enabled.
- **`attribute_status`**: the latest values and classifications for all metrics in one row per tenant and car.
- **`health_events`**: warning or critical intervals for a tenant, car, and metric.

The application tables use tenant-aware composite keys. Both `attribute_status` and `health_events` reference `(tenant_id, vin)` in `cars` and cascade on car deletion.

## Tenancy and operational scope

- All registered HTTP endpoints are authenticated by API-key middleware.
- Tenant ID and a tenant-enriched logger flow through request context.
- Repository data queries require the tenant context and scope reads and writes by tenant ID.
- `cmd/provision` is an out-of-band CLI that stores a key hash and prints the generated API key once.
- The checked-in schema is a fresh-database POC schema, not a migration of existing data.
- MySQL HTTP isolation tests run only when `MYSQL_TEST_DSN` is set; `docker-compose.yml` provides a disposable local database.
