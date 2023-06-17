FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o email-service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/email-service .
CMD ["./email-service"]