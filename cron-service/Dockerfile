# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache tzdata

WORKDIR /app

COPY api-proto ./api-proto

WORKDIR /app/cron-service

COPY cron-service/go.mod cron-service/go.sum ./
RUN go mod download
COPY cron-service .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o app ./cmd

# ---------- Runtime stage ----------
FROM scratch

WORKDIR /app

COPY --from=builder /app/cron-service/app /app/app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Moscow

EXPOSE 8081

ENTRYPOINT ["/app/app"]
