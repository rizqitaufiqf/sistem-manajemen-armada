# Base builder image
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy all source code
COPY . .

# Build API service
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

# Build MQTT Subscriber service
RUN CGO_ENABLED=0 GOOS=linux go build -o /mqtt-subscriber ./cmd/mqtt-subscriber

# Build RabbitMQ Worker service
RUN CGO_ENABLED=0 GOOS=linux go build -o /rabbitmq-worker ./cmd/rabbitmq-worker

# Build MQTT Publisher script
RUN CGO_ENABLED=0 GOOS=linux go build -o /mqtt-publisher ./cmd/mqtt-publisher


# Final stage: production-ready slim image
FROM alpine:latest

WORKDIR /app

# Copy binaries from the builder stage
COPY --from=builder /api /app/api
COPY --from=builder /mqtt-subscriber /app/mqtt-subscriber
COPY --from=builder /rabbitmq-worker /app/rabbitmq-worker
COPY --from=builder /mqtt-publisher /app/mqtt-publisher

# Expose ports
EXPOSE 8080