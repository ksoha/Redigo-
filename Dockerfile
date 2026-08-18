# --- Build stage ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies separately from source for faster rebuilds
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /redigo ./cmd/redigo

# --- Run stage ---
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /redigo /app/redigo

# AOF file lives here — mounted as a volume so data survives container restarts
VOLUME ["/app/data"]

EXPOSE 6380

ENTRYPOINT ["/app/redigo"]