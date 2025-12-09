FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go build -o app ./cmd

EXPOSE 8080

CMD ["./app"]
