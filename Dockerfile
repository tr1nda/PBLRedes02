# Etapa de build
FROM golang:1.22-alpine AS builder

ARG TARGET=client
RUN echo "Building for $TARGET"
WORKDIR /app

# Copia go.mod primeiro (útil para cache), mas agora também o código completo
COPY go.mod ./
COPY . ./
COPY data/ ./data/

# Executa o tidy agora, pois precisa analisar os imports do código
RUN go mod tidy

# Compila o binário
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$TARGET ./$TARGET/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app ./${TARGET}/main.go

# Imagem final
FROM alpine:latest

ARG TARGET=client
RUN echo "Building for $TARGET"
WORKDIR /root/

# Copia o binário da etapa de build
COPY --from=builder /app/data /root/data
COPY --from=builder /app/bin/app /usr/local/bin/app

CMD ["/usr/local/bin/app"]