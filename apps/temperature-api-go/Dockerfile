# syntax=docker/dockerfile:1

# --- Build stage -------------------------------------------------------------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies first.
COPY go.mod go.sum ./
RUN go mod download

# Build a static binary so we can use a minimal runtime image.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o temperature-api-go .

# --- Runtime stage -----------------------------------------------------------
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata \
    && adduser -D -u 10001 appuser

WORKDIR /app
COPY --from=builder /app/temperature-api-go .

USER appuser

EXPOSE 8081
ENV PORT=:8081

CMD ["./temperature-api-go"]
