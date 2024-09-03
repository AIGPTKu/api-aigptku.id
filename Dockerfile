# Stage 1: Build Stage
FROM golang:1.20.7-alpine AS builder

RUN apk update && apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary cmd/*.go

# Stage 2: Move binary file
FROM alpine

WORKDIR /app

COPY --from=builder /app/binary .

CMD ["/app/binary"]