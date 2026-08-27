# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
