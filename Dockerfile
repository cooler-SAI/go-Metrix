FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]