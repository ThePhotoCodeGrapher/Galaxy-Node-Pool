# Dockerfile for galaxy-node-pool
FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o pool-server ./cmd/pool-server

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/pool-server .
COPY configs/example.yaml ./configs/example.yaml
EXPOSE 50051
ENTRYPOINT ["./pool-server", "--config", "./configs/example.yaml"]
