FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go build -o users-service ./cmd/users

EXPOSE 8080

CMD ["./users-service"]

