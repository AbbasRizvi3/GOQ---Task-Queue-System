FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/internal/static ./internal/static
EXPOSE 8000
CMD ["./main"]