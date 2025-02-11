FROM golang:1.23.5 AS builder


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR /app/cmd/mcu_service

RUN go build -o mcu_service

FROM debian:bookworm-slim 

WORKDIR /app

COPY --from=builder /app/cmd/mcu_service .
COPY --from=builder /app/config .

# Add this to your Dockerfile to install wait-for-it
RUN apt-get update && apt-get install -y wait-for-it

# Add this to your entrypoint to wait for Kafka to be ready before starting mcu_service
CMD ["wait-for-it", "kafka:9092", "--", "./mcu_service"]
