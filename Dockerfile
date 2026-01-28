# syntax=docker/dockerfile:1.6

# --- Stage 1: Builder ---
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /src

# Cache deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source
COPY . .

# Build
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/go-boilerplate .

# --- Stage 2: Runtime ---
FROM scratch

# CA certs (for TLS to Supabase, etc.)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Timezone data (optional, but helpful for logs)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# App binary
COPY --from=builder /out/go-boilerplate /go-boilerplate

# Create writable dirs (scratch has nothing, so we "copy" empty dirs from builder)
# We'll run as non-root numeric user (no passwd needed)
COPY --from=builder /tmp /tmp

# Runtime env: avoid "/.cache" permission errors
ENV HOME=/tmp \
    XDG_CACHE_HOME=/tmp/.cache \
    GOMODCACHE=/tmp/gomodcache \
    TZ=UTC

EXPOSE 5000

# Run as unprivileged user (numeric works without /etc/passwd)
USER 65534:65534

ENTRYPOINT ["/go-boilerplate"]
