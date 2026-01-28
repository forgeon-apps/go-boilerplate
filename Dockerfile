# --- Stage 1: Builder ---
FROM golang:1.23-alpine AS builder

# Create unprivileged user
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install certs
RUN apk add --no-cache ca-certificates

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary with optimizations
# -ldflags="-s -w" removes debug info to shrink the size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o go-boilerplate .

# --- Stage 2: Final (Production) ---
FROM scratch AS final

# Standard practice: Use uppercase for AS to avoid warnings
COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /build/go-boilerplate /go-boilerplate

# OPTIONAL: Only uncomment if you strictly need .env baked in 
# and ensure it exists in your repo.
# COPY --from=builder /build/.env / .env

EXPOSE 5000

USER nobody:nobody

ENTRYPOINT ["/go-boilerplate"]