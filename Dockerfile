# syntax=docker/dockerfile:1.6

# ----------------------------
# Stage 1: Builder
# ----------------------------
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata git

WORKDIR /src

# Cache deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source
COPY . .

# Build (supports buildx arch)
ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/go-boilerplate .

# ----------------------------
# Stage 2: Runtime (nonroot, sane /tmp)
# ----------------------------
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

# TLS certs for HTTPS (Supabase, etc.)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Timezone data (optional; remove if you want ultra-minimal)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /out/go-boilerplate /go-boilerplate

# Writable cache/home (prevents "/tmp/.cache permission denied")
ENV HOME=/tmp \
    XDG_CACHE_HOME=/tmp/.cache \
    TZ=UTC

EXPOSE 5000

ENTRYPOINT ["/go-boilerplate"]
