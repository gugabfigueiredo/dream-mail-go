FROM golang:1.20 AS builder
# Build the application
ADD . /go/src/github/gugabfigueiredo/dream-mail-go
WORKDIR /go/src/github/gugabfigueiredo/dream-mail-go/cmd/dream-mail-server
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ../../dream-mail-go .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github/gugabfigueiredo/dream-mail-go/dream-mail-go .
COPY .env .
CMD ["sh", "-c", "source .env && ./dream-mail-go"]