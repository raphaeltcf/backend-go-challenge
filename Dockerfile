FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/app .
CMD ["./app"]