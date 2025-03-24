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
RUN go build -o main code_processor/main.go

# Финальный этап (для уменьшения размера образа)
FROM alpine:latest AS runner

# Устанавливаем рабочую директорию
WORKDIR /code_processor

# Копируем собранное приложение из этапа сборки
COPY --from=build /build/main ./main
COPY --from=build /build/code_processor/config/config.yml ./config.yml

# Даем права на выполнение
RUN chmod +x /code_processor/main

# Указываем команду для запуска
CMD ["/code_processor/main", "--config=/code_processor/config.yml"]