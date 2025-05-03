# Etapa de build
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o pbl-redes

# Etapa final
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pbl-redes .

CMD ["./pbl-redes"]
