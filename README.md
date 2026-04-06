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

### Вариант 1: Использование Makefile (рекомендуется)

```bash
# Показать все доступные команды
make help

# Полная настройка проекта (установка зависимостей, запуск инфраструктуры, миграции)
make setup

# Или пошагово:
# 1. Запустить PostgreSQL и Redis
make up-infra

# 2. Применить миграции
make migrate-up

# 3. Заполнить тестовыми данными (опционально)
make db-seed

# 4. Запустить все сервисы в Docker
make up

# Посмотреть логи
make logs

# Посмотреть статус контейнеров
make ps
```

### Вариант 2: Ручная настройка

#### 1. Клонировать репозиторий

```bash
git clone <repo-url>
cd dns-proxy-pet
```

#### 2. Настроить переменные окружения

```bash
cp .env.example .env
# Отредактировать .env файл
```

#### 3. Запустить PostgreSQL и Redis

```bash
docker compose up -d postgres redis
```

#### 4. Применить миграции

```bash
./scripts/migrate.sh
```

#### 5. Запустить сервисы локально

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

### Makefile команды

Проект использует Makefile для упрощения разработки. Основные команды:

#### Управление инфраструктурой
```bash
make up-infra          # Запустить PostgreSQL и Redis
make down              # Остановить все сервисы
make restart-infra     # Перезапустить инфраструктуру
make ps                # Показать статус контейнеров
make logs              # Показать логи всех сервисов
make logs-db           # Показать логи PostgreSQL
make logs-redis        # Показать логи Redis
```

#### Управление базой данных
```bash
make migrate-up        # Применить все миграции
make migrate-down      # Откатить последнюю миграцию
make migrate-reset     # Сбросить БД и применить миграции заново
make db-tables         # Показать все таблицы
make db-shell          # Открыть PostgreSQL shell
make db-seed           # Заполнить БД тестовыми данными
make db-clean          # Удалить все таблицы
```

#### Redis
```bash
make redis-shell       # Открыть Redis CLI
make redis-flush       # Очистить все данные Redis
```

#### Разработка
```bash
make run-dns           # Запустить DNS сервер локально
make run-proxy         # Запустить Proxy сервер локально
make run-api           # Запустить API сервер локально
make run-bot           # Запустить Telegram бот локально
```

#### Сборка
```bash
make build-all         # Собрать все бинарники
make build-dns         # Собрать DNS сервер
make build-proxy       # Собрать Proxy сервер
make build-api         # Собрать API сервер
make build-bot         # Собрать Telegram бот
```

#### Тестирование
```bash
make test              # Запустить все тесты
make test-coverage     # Запустить тесты с coverage
make lint              # Запустить линтер
make fmt               # Форматировать код
make vet               # Запустить go vet
```

#### Очистка
```bash
make clean             # Удалить контейнеры, volumes, бинарники
make clean-cache       # Очистить Go кеш
```

#### Информация
```bash
make help              # Показать все команды
make info              # Показать информацию о проекте
```

### Установка зависимостей

```bash
make tidy
# или
go mod download
```

### Запуск тестов

```bash
make test
# или
go test ./...
```

### Сборка

```bash
make build-all
# или вручную
go build -o bin/dns-server ./cmd/dns-server
go build -o bin/proxy-server ./cmd/proxy-server
go build -o bin/api-server ./cmd/api-server
go build -o bin/telegram-bot ./cmd/telegram-bot
```

### Docker

```bash
# Собрать все образы
make build
# или
docker compose build

# Запустить все сервисы
make up
# или
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

### Полезные make команды

- make setup          # Полная автоматическая настройка
- make up-infra       # Запуск PostgreSQL + Redis
- make migrate-up     # Применить миграции
- make db-seed        # Тестовые данные
- make db-clean       # Очистить БД
- make ps             # Статус контейнеров
- make help           # Все команды