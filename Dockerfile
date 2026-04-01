# Stage 1: build
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for fetching deps and TLS)
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy dependency manifests first to leverage Docker layer cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm \
    go build -trimpath \
    -ldflags="-s -w" \
    -o sql_exporter \
    ./cmd/sql_exporter/

# Stage 2: minimal runtime image
FROM gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.title="mssql_exporter" \
      org.opencontainers.image.description="Prometheus SQL Exporter for MSSQL and other databases" \
      org.opencontainers.image.source="https://github.com/mateusgsan/mssql_exporter" \
      org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /build/sql_exporter /bin/sql_exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 9399

USER nonroot:nonroot

ENTRYPOINT ["/bin/sql_exporter"]
