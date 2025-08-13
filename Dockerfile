# --- build stage ---
FROM golang:1.24.5 AS build
ENV GOTOOLCHAIN=local
WORKDIR /app

# Сначала зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем всё остальное
COPY . .

# Собираем бинарь (у тебя main.go лежит в ./cmd/main)
RUN go build -o url-shorter ./cmd/main

# --- run stage ---
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# Копируем бинарь из build stage
COPY --from=build /app/url-shorter /app/url-shorter

# Переменная окружения — путь к конфигу
ENV CONFIG_PATH=/config/prod.yaml

# Порт, который слушает твой сервис
EXPOSE 8080

# Запуск приложения
ENTRYPOINT ["/app/url-shorter"]