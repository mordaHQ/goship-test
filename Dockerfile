# ЭТАП 1: Сборщик (Builder)
# Берем тяжелый образ с Go, чтобы скомпилировать
FROM golang:alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o server main.go

# ЭТАП 2: Финальный образ (Runner)
# Берем крошечный Alpine Linux (без Go, без мусора)
FROM alpine:latest
WORKDIR /root/

# Копируем ТОЛЬКО готовый файл из первого этапа
COPY --from=builder /app/server .

# Запускаем
CMD ["./server"]
