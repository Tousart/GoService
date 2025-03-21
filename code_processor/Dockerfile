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
RUN go build -o main .

# Финальный этап (для уменьшения размера образа)
FROM alpine:latest AS runner

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение из этапа сборки
COPY --from=build /build/main ./main
COPY --from=build /build/config/config.yml ./config.yml

# Даем права на выполнение
RUN chmod +x /app/main

# Указываем команду для запуска
CMD ["/app/main", "--config=/app/config.yml"]