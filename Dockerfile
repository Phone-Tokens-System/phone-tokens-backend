FROM golang:1.25-alpine AS BUILDER

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd

FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./main"]
