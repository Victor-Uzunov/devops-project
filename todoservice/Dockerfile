FROM golang:1.23.2 AS builder

WORKDIR /app

COPY . .
COPY .env ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o bin/main ./cmd/todoapp/main.go


FROM alpine:3
WORKDIR /
COPY --from=builder /app/bin/main /main
COPY --from=builder /app/.env ./

ENTRYPOINT ["./main"]
EXPOSE 5000

