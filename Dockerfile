# ---------- Stage 1: Build ----------
FROM golang:alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download || (sleep 2 && go mod download)

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api

# ---------- Stage 2: Runtime ----------
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder --chown=appuser:appgroup /app/api .
COPY --from=builder --chown=appuser:appgroup /app/migrations ./migrations

USER appuser

EXPOSE 8080

CMD ["./api"]