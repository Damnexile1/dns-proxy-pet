# DNS + Proxy Service

Коммерческий сервис обхода блокировок через DNS + HTTP/HTTPS прокси с системой подписок через Telegram бота.

## Структура проекта

```
.
├── cmd/                    # Точки входа для сервисов
│   ├── dns-server/        # DNS сервер
│   ├── proxy-server/      # HTTP/HTTPS прокси
│   ├── api-server/        # REST API
│   └── telegram-bot/      # Telegram бот
├── internal/              # Внутренняя бизнес-логика
│   ├── config/           # Конфигурация
│   ├── database/         # Подключение к БД
│   ├── models/           # Модели данных
│   ├── repository/       # Слой работы с БД
│   ├── service/          # Бизнес-логика
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # Middleware
│   └── utils/            # Утилиты
├── pkg/                   # Публичные пакеты
│   ├── logger/           # Логирование
│   ├── cache/            # Кеширование
│   └── auth/             # Аутентификация
├── migrations/            # SQL миграции
├── scripts/              # Скрипты
├── docs/                 # Документация
├── config.yaml           # Конфигурация
├── .env                  # Переменные окружения
└── docker-compose.yml    # Docker Compose
```

## Требования

- Go 1.22+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone <repo-url>
cd dns-proxy-pet
```

### 2. Настроить переменные окружения

```bash
cp .env.example .env
# Отредактировать .env файл
```

### 3. Запустить PostgreSQL и Redis

```bash
docker compose up -d postgres redis
```

### 4. Применить миграции

```bash
# Установить golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Применить миграции
migrate -path migrations -database "postgresql://dnsproxy:dev_password_123@localhost:5432/dnsproxy?sslmode=disable" up
```

### 5. Запустить сервисы

```bash
# DNS сервер
go run cmd/dns-server/main.go

# Proxy сервер
go run cmd/proxy-server/main.go

# API сервер
go run cmd/api-server/main.go

# Telegram бот
go run cmd/telegram-bot/main.go
```

## Разработка

### Установка зависимостей

```bash
go mod download
```

### Запуск тестов

```bash
go test ./...
```

### Сборка

```bash
# Собрать все сервисы
go build -o bin/dns-server ./cmd/dns-server
go build -o bin/proxy-server ./cmd/proxy-server
go build -o bin/api-server ./cmd/api-server
go build -o bin/telegram-bot ./cmd/telegram-bot
```

### Docker

```bash
# Собрать все образы
docker compose build

# Запустить все сервисы
docker compose up -d

# Посмотреть логи
docker compose logs -f

# Остановить все сервисы
docker compose down
```

## Конфигурация

Конфигурация загружается из `config.yaml` и может быть переопределена через переменные окружения из `.env`.

Основные параметры:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL
- `REDIS_HOST`, `REDIS_PORT` - Redis
- `DNS_PORT` - порт DNS сервера (по умолчанию 53)
- `PROXY_HTTP_PORT`, `PROXY_HTTPS_PORT` - порты прокси
- `API_PORT` - порт API сервера
- `TELEGRAM_BOT_TOKEN` - токен Telegram бота
- `JWT_SECRET` - секрет для JWT токенов

## База данных

### Таблицы

- `users` - пользователи
- `subscriptions` - подписки
- `proxy_credentials` - учетные данные для прокси
- `payments` - платежи
- `traffic_usage` - использование трафика
- `blocked_domains` - заблокированные домены

### Миграции

Миграции находятся в папке `migrations/` и применяются с помощью `golang-migrate`.

```bash
# Применить все миграции
migrate -path migrations -database "postgresql://..." up

# Откатить последнюю миграцию
migrate -path migrations -database "postgresql://..." down 1

# Создать новую миграцию
migrate create -ext sql -dir migrations -seq <migration_name>
```

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - регистрация
- `POST /api/v1/auth/verify` - проверка токена

### Пользователь
- `GET /api/v1/user/profile` - профиль
- `GET /api/v1/user/subscription` - статус подписки
- `GET /api/v1/user/credentials` - данные для подключения
- `GET /api/v1/user/traffic` - статистика трафика

### Подписки
- `POST /api/v1/subscription/create` - создать подписку
- `POST /api/v1/subscription/renew` - продлить
- `POST /api/v1/subscription/cancel` - отменить автопродление

### Платежи
- `POST /api/v1/payment/create` - создать платеж
- `POST /api/v1/payment/webhook` - webhook от платежной системы
- `GET /api/v1/payment/history` - история платежей

## Telegram Bot

### Команды

- `/start` - начало работы, регистрация
- `/status` - статус подписки
- `/pay` - оплатить подписку
- `/settings` - получить настройки подключения
- `/help` - инструкции по настройке
- `/support` - связаться с поддержкой
- `/cancel` - отменить автопродление

## Лицензия

Proprietary
