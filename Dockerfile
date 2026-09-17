FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/expiration ./cmd/expiration


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/bin/api ./api

EXPOSE 8080

CMD ["./api"]

FROM alpine:3.22 AS expiration-worker

WORKDIR /app

COPY --from=builder /bin/expiration-worker ./expiration-worker

CMD ["./expiration-worker"]