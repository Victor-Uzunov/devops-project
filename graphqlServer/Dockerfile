FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY .env ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bin/main ./cmd/main.go

FROM alpine:3
WORKDIR /root/
COPY --from=builder /app/bin/main .
COPY --from=builder /app/.env ./
CMD ["./main"]
