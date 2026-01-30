FROM golang:alpine AS builder
WORKDIR /app

# Нужно для работы SQLite на Alpine
ENV CGO_ENABLED=0 

# Инициализируем модуль и качаем зависимости
COPY main.go .
RUN go mod init goship-app
RUN go get modernc.org/sqlite
RUN go mod tidy

# Компилируем
RUN go build -o server main.go

FROM alpine:latest
WORKDIR /root/
# Копируем сервер и html
COPY --from=builder /app/server .
COPY index.html .

CMD ["./server"]
