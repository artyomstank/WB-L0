# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o wb-service ./cmd/main.go

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/wb-service .
COPY ./web ./web
COPY .env .env

EXPOSE 8081
CMD ["./wb-service"]