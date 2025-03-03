# Stage 1
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go build -a -installsuffix cgo -o main ./cmd

# Stage 2
FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 80

CMD ["./main"]