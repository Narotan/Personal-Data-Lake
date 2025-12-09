FROM golang:1.25.1-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o data-lake cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget postgresql-client

# Создаем непривилегированного пользователя
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

WORKDIR /app

COPY --from=builder /build/data-lake .
COPY db/migrations /app/migrations

# Меняем владельца файлов на appuser
RUN chown -R appuser:appgroup /app

# Переключаемся на непривилегированного пользователя
USER appuser

EXPOSE 8080

CMD ["./data-lake"]
р