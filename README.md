# Calculator Platform

A small full-stack calculator with a React interface and a Go HTTP API. It supports addition, subtraction, multiplication, division, powers, square roots, and percentages.

## Architecture

The repository has two deployable build inputs and one runtime process:

- `internal/calculator` keeps calculation rules, request validation, HTTP handling, and response mapping in one feature package with separate responsibilities.
- `cmd/api` separates process startup from server composition and serves the built frontend in production.
- `frontend/src/pages` contains the thin page composition entry.
- `frontend/src/features/calculator` owns the calculator model, API boundary, workflow hook, components, and focused tests.

The Go service uses only the standard library. The feature stays grouped by behavior, while small files separate transport parsing, validation, result mapping, and calculation rules. The React page delegates interaction state and the calculation workflow to a feature-local hook. This keeps responsibilities clear without adding interfaces or shared abstractions that have only one implementation.

Production uses a multi-stage container. Node builds the static assets, Go builds one static binary, and the final image contains only those artifacts. The browser calls the API on the same origin, so no CORS configuration is needed.

## API

`POST /api/calculate` accepts JSON. Binary operations use both operands:

```json
{
  "operation": "multiply",
  "left": 7,
  "right": 5
}
```

Square root accepts only `left`:

```json
{
  "operation": "square_root",
  "left": 81
}
```

Percentage calculates `left` percent of `right`. For example, 15 percent of 200 returns 30.

Successful responses contain a result:

```json
{
  "result": 35
}
```

Validation errors use stable codes and readable messages:

```json
{
  "error": {
    "code": "division_by_zero",
    "message": "Cannot divide by zero."
  }
}
```

The API uses Go `float64` values and rejects results that are not finite. `GET /healthz` provides a health check.

## Local development

Requirements:

- Go 1.26
- Node.js 24
- Docker for the production build

Run the API:

```sh
go run ./cmd/api
```

In another terminal, run the frontend development server:

```sh
cd frontend
npm ci
npm run dev
```

Vite proxies `/api` and `/healthz` to `http://localhost:8080`.

## Verification

```sh
gofmt -w .
go vet ./...
go test -race -coverprofile=coverage.out ./...
go test -bench=. -benchmem ./internal/calculator

cd frontend
npm run lint
npm run test:coverage
npm run build
```

Build and run the same container used in deployment:

```sh
docker build -t calculator-platform .
docker run --rm -p 8080:8080 calculator-platform
```

Open `http://localhost:8080`. The GitHub Actions workflow repeats formatting, static analysis, tests, coverage, builds, and container smoke tests on every change.

## Deployment

`render.yaml` defines one Docker web service with health checks and deploys only after CI passes.

To provision it, create a new Blueprint in the Render dashboard, connect this repository, and apply the detected configuration. Render then builds the `Dockerfile` and monitors `/healthz`.

The free service plan can sleep between requests, so the first request after an idle period may take longer. Change the plan before launch if consistent response time is required.
