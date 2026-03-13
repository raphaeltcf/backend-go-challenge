FROM golang:1.26 AS tester
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go test ./...

FROM golang:1.26 AS builder
WORKDIR /app
COPY --from=tester /app .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]