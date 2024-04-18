FROM golang:1.22-bookworm AS builder

# Set Go env
ENV GOOS=linux CGO_ENABLED=0

WORKDIR /workspace

# Build Go binary
COPY . .
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /workspace/goac

# Deployment container
FROM gcr.io/distroless/static-debian12

COPY --from=builder /workspace/goac /goac
ENTRYPOINT ["/goac"]
