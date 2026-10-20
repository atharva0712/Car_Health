# Car Health

A Go and MySQL vehicle-health monitoring service. The backend registers cars, accepts generated telemetry, tracks the latest metric status, and records warning and critical health events.

## Current features

- Register, list, retrieve, delete, and toggle telemetry simulation for cars.
- Receive telemetry for tire pressure, engine temperature, battery voltage, oil level, and brake condition.
- Classify telemetry using thresholds in `config.yaml`.
- Read a car's current status and health events.
- Run a separate telemetry generator that discovers cars from the backend and submits generated telemetry.

## Architecture

The backend follows a handler → service → repository layout:

- `cmd/backend` wires the HTTP server and application dependencies.
- `internal/handler` implements HTTP endpoints.
- `internal/service` processes car operations and telemetry/status logic.
- `internal/repository` performs MySQL queries.
- `internal/db/migrations/schema.sql` defines the current database schema.
- `cmd/generator` runs the telemetry generator in a separate process.
- `config.yaml` defines metric ranges, generation distribution, interval, and backend URL.

See [docs/architecture.md](docs/architecture.md) for the request and data flows, and [docs/multi-tenancy.md](docs/multi-tenancy.md) for the implemented POC tenancy model and its limitations.

## Requirements

- Go version from `go.mod` (currently Go 1.25.6).
- MySQL.

## Local setup

1. Create a local `.env` from `.env.example` and set the database credentials.
2. Apply `internal/db/migrations/schema.sql` to a **fresh database**. The POC schema is not a migration and must not be applied to a database containing existing data.
3. Provision a tenant from a trusted terminal:

   ```powershell
   go run ./cmd/provision --name "Development"
   ```

   Save the printed API key securely; it is shown only once.
4. Set that tenant's key in the environment used by the telemetry generator:

   ```powershell
   $env:API_KEY = "paste-the-provisioned-key"
   ```

5. Start the backend:

   ```powershell
   go run ./cmd/backend
   ```

6. In another terminal, set `API_KEY` to the same tenant key and start the telemetry generator:

   ```powershell
   go run ./cmd/generator
   ```

The backend listens on `BACKEND_PORT` (default `8080`). The generator uses the backend URL in `config.yaml`. Every API endpoint requires the `X-API-Key` header.

## Multi-tenancy

- Tenants share one MySQL schema; each car, status, and event is scoped by `tenant_id`.
- `X-API-Key` is hashed and resolved by middleware; tenant identity is never accepted from request data.
- Context flows from middleware through handlers and services to tenant-scoped repository queries.
- Composite foreign keys enforce that status and event rows belong to a car in the same tenant.
- This POC prioritizes simple isolation; database-per-tenant offers stronger separation at higher operational cost.

The schema in this POC is for fresh databases only. It does not migrate existing data; do not use it to replace a populated database without first writing a backfill migration.

## Tests

Run unit tests:

```powershell
go test ./...
```

For MySQL integration tests, start the disposable database with `docker compose up -d mysql`, set `MYSQL_TEST_DSN` to the service's connection string, then run:

```powershell
$env:MYSQL_TEST_DSN = "root:car_health_test@tcp(127.0.0.1:3307)/car_monitor?parseTime=true"
go test ./...
```

## HTTP API

| Method | Path | Purpose |
| --- | --- | --- |
| `POST` | `/cars` | Register a car |
| `GET` | `/cars` | List cars |
| `GET` | `/cars/{vin}` | Retrieve one car |
| `DELETE` | `/cars/{vin}` | Delete a car |
| `PATCH` | `/cars/{vin}/simulation` | Enable or disable simulation |
| `POST` | `/telemetry` | Submit telemetry |
| `GET` | `/cars/{vin}/status` | Retrieve all metric statuses |
| `GET` | `/cars/{vin}/status?attribute=tire_pressure` | Retrieve one metric status |
| `GET` | `/cars/{vin}/events?date=YYYY-MM-DD` | Retrieve health events |

Every API route requires a valid tenant API key in the `X-API-Key` header.
