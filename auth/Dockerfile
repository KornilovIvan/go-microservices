FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./bin/grpc_server ./cmd/grpc_server

FROM alpine:3.20

WORKDIR /root

COPY --from=builder /build/bin/grpc_server .
COPY prod.env .

CMD ["./grpc_server", "--config-path", "prod.env"]
