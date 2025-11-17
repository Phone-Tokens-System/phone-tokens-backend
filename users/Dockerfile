FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o users-service ./cmd/users

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/users-service .

ENV HTTP_PORT=8080

EXPOSE 8080

CMD ["./users-service"]

