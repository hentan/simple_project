# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект, чтобы были cmd/, internal/, и т.д.
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/app

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates

RUN apk add --no-cache curl && \
    curl -L https://github.com/pressly/goose/releases/download/v3.21.1/goose_linux_x86_64 -o /usr/local/bin/goose && \
    chmod +x /usr/local/bin/goose && \
    apk del curl

COPY --from=builder /out/app ./app
COPY internal/migrations ./migrations
COPY .env ./.env

EXPOSE 8080
CMD ["./app"]
