# Etapa de build
FROM golang:1.24 AS builder

ARG TARGET

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o main "./cmd/${TARGET}"

# Etapa final
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
