# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY email-sender/go.mod email-sender/go.sum ./
RUN go mod download

COPY email-sender .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o app ./cmd

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/app /app/app

ENTRYPOINT ["/app/app"]
