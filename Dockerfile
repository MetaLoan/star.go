# Multi-stage build for Star API (v2) with Swiss Ephemeris support

# ========== Builder ==========
FROM golang:1.21-bookworm AS builder

WORKDIR /app

# Swiss Ephemeris headers/libs for CGO build
RUN apt-get update && apt-get install -y --no-install-recommends \
    libswe-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Go dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source
COPY backend .

# Build with Swiss Ephemeris
ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-lm"
RUN go build -tags swe -o /out/star-api .

# ========== Runtime ==========
FROM debian:bookworm-slim AS runtime

WORKDIR /app

# Runtime deps
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libswe2.0 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/star-api /app/star-api

EXPOSE 8080

CMD ["./star-api"]
