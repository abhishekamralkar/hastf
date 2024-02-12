# Stage 1: Build the binary
FROM golang:1.22.0 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o hastf .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/hastf .

CMD ["./hastf"]