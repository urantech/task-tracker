# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY api-proto ./api-proto

WORKDIR /app/backend-api

COPY backend-api/go.mod backend-api/go.sum ./
RUN go mod download

COPY backend-api .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o app ./cmd

# ---------- Runtime stage ----------
FROM scratch

WORKDIR /app
COPY --from=builder /app/backend-api/app /app/app

EXPOSE 8080
ENTRYPOINT ["/app/app"]
