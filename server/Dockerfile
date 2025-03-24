# Используем образ Go с нужной версией
FROM golang:1.23.5-alpine AS build

# Устанавливаем рабочую директорию
WORKDIR /build

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main server/main.go

# Финальный этап (для уменьшения размера образа)
FROM alpine:latest AS runner

# Устанавливаем рабочую директорию
WORKDIR /http_server

# Копируем собранное приложение из этапа сборки
COPY --from=build /build/main ./main
COPY --from=build /build/server/config/config.yml ./config.yml

# Устанавливаем curl
RUN apk add --no-cache curl

# Даем права на выполнение
RUN chmod +x /http_server/main

# Указываем команду для запуска
CMD ["/http_server/main", "--config=/http_server/config.yml"]