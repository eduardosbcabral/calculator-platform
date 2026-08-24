# Calculator Platform

A small full-stack calculator with a React interface and a Go HTTP API. It supports addition, subtraction, multiplication, division, powers, square roots, and percentages.

## Architecture

The repository has two deployable build inputs and one runtime process:

- `internal/calculator` keeps calculation rules, request validation, HTTP handling, and response mapping in one feature package with separate responsibilities.
- `cmd/api` separates process startup from server composition and serves the built frontend in production.
- `frontend/src/pages` contains the thin page composition entry.
- `frontend/src/features/calculator` owns the calculator model, API boundary, workflow hook, components, and focused tests.

The Go service uses only the standard library. The feature stays grouped by behavior, while small files separate transport parsing, validation, result mapping, and calculation rules. The React page delegates interaction state and the calculation workflow to a feature-local hook. This keeps responsibilities clear without adding interfaces or shared abstractions that have only one implementation.

Production uses a multi-stage container. Node builds the static assets, Go builds one static binary, and a small Alpine runtime runs the non-root process and serves those artifacts. The browser calls the API on the same origin, so no CORS configuration is needed.

## API

`POST /api/calculate` accepts JSON. The supported operations are:

| Operation | Calculation |
| --- | --- |
| `add` | `left + right` |
| `subtract` | `left - right` |
| `multiply` | `left * right` |
| `divide` | `left / right` |
| `power` | `left` raised to `right` |
| `percentage` | `left` percent of `right` |
| `square_root` | Square root of `left` |

Binary operations use both operands:

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

Requests must use `application/json`, contain one object with known fields, and stay within 4 KiB. The API uses Go `float64` values and rejects invalid operands and results that are not finite. `GET /healthz` provides a health check.

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

## Testing and coverage

The backend test suite covers every calculator operation, arithmetic edge cases, request validation, JSON errors, HTTP routing, health checks, and static frontend serving. It also includes a benchmark for the calculation path.

The frontend has seven Vitest and React Testing Library tests covering binary and unary calculations, request payloads, input validation, API errors, loading state, and result rendering. API calls are mocked so the frontend behavior is tested independently from the backend.

Current coverage on `main`:

| Scope | Coverage |
| --- | ---: |
| Go, all packages | 81.2% statements |
| Go, `internal/calculator` | 98.7% statements |
| React statements | 87.50% |
| React branches | 85.71% |
| React functions | 93.75% |
| React lines | 86.79% |

Frontend coverage has an enforced 80% minimum for statements, branches, functions, and lines. The [GitHub Actions workflow](https://github.com/eduardosbcabral/calculator-platform/actions/workflows/ci.yml) runs both suites and uploads the Go coverage profile and frontend HTML and JSON reports in an artifact named `coverage`.

## Verification

```sh
gofmt -w .
go vet ./...
go test -race -coverprofile=coverage.out ./...
go test -bench=. -benchmem ./internal/calculator

cd frontend
npm ci
npm run lint
npm run test:coverage
npm run build
```

Build and run the same container used in deployment:

```sh
docker build -t calculator-platform .
docker run --rm -p 8080:8080 calculator-platform
```

Open `http://localhost:8080`. CI repeats formatting, static analysis, tests, coverage, builds, and container smoke tests on every change.

## Deployment

The application runs as one Docker service on Coolify. The platform builds the root `Dockerfile`, exposes port 8080, and monitors `/healthz`.

Live application: [calculator.propeller.com.br](https://calculator.propeller.com.br)
