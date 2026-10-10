# temperature-api-go

A minimal Go service that exposes temperature sensor data from the **same
PostgreSQL database** used by the `smart_home` service. It is a port of the
`temperature-api-net` service.

- Standard library `net/http` only (no web framework).
- Single external dependency: `github.com/jackc/pgx/v5` (raw SQL, no ORM).
- Small Alpine-based Docker image.
- Returns the exact JSON shape of `services.TemperatureResponse`.

## Endpoints

| Method | Route                     | Go caller                                         |
| ------ | ------------------------- | ------------------------------------------------- |
| GET    | `/temperature?location=…` | `TemperatureService.GetTemperature(location)`     |
| GET    | `/temperature/{sensorID}` | `TemperatureService.GetTemperatureByID(sensorID)` |
| GET    | `/health`                 | Container health check                            |
| GET    | `/`                       | Basic service info                                |

Example response body (matches `TemperatureResponse`):

```json
{
  "value": 22.5,
  "unit": "C",
  "timestamp": "2026-10-10T12:00:00Z",
  "location": "Living Room",
  "status": "active",
  "sensor_id": "1",
  "sensor_type": "temperature",
  "description": "Living Room thermometer"
}
```

## Configuration

| Variable       | Default                                                     | Description                          |
| -------------- | ----------------------------------------------------------- | ------------------------------------ |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/smarthome`      | PostgreSQL connection string (URI or key/value) |
| `PORT`         | `:8080`                                                      | Address/port the HTTP server binds   |

## Run with Docker Compose

From the `apps` directory:

```bash
docker compose up --build
```

The service listens on container port `8080` and is published on host port `8081`.

## Run locally

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5432/smarthome" \
  go run .
```

## Build

```bash
go build -o temperature-api-go .
docker build -t temperature-api-go .
```
