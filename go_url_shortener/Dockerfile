# Build
FROM golang:1.23-alpine AS build

WORKDIR /app

# Copia go.mod e go.sum primeiro para cache de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código
COPY . .

# Build binário
RUN go build -o /url_shortener ./cmd/main.go

# Runtime
FROM alpine:3.18

WORKDIR /app

COPY --from=build /url_shortener .
COPY .env .

EXPOSE 8080

CMD ["./url_shortener"]
