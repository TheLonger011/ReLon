FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /relon ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /relon /app/relon
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/relon"]