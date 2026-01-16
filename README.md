# Task Tracker

Микросервисное приложение для управления задачами. Система включает в себя REST и gRPC API для работы с задачами и пользователями, cron-сервис для автоматических операций и сервис отправки email-уведомлений.

## Стек технологий

- **Язык**: Go 1.25.5
- **Базы данных**: PostgreSQL
- **Очереди сообщений**: Apache Kafka
- **Контейнеризация**: Docker, Docker Compose
- **HTTP-фреймворк**: Chi
- **gRPC**: Google Protocol Buffers
- **Аутентификация**: JWT
- **Миграции**: Goose
- **Тестирование**: Testcontainers

## Предварительные требования

- Go 1.25 или выше
- Docker и Docker Compose
- Git

## Установка и настройка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/urantech/task-tracker-go.git
   cd task-tracker
   ```

2. Скопируйте файл с примером переменных окружения и настройте его:
   ```bash
   cp .env.example .env
   ```
   Отредактируйте `.env` файл, заполнив необходимые значения (пароли БД, JWT секрет и т.д.).

3. Установите зависимости и запустите миграции (автоматически при запуске сервисов через Docker).

## Использование

Проект использует [Task](https://github.com/go-task/task) для автоматизации команд. Установите его, если необходимо:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

### Доступные команды (через Taskfile)

#### Управление инфраструктурой
| Команда | Описание |
|---------|----------|
| `task start-infra` | Запускает инфраструктуру: PostgreSQL, Kafka и Mailhog (SMTP-сервер для тестирования) |
| `task stop-all` | Останавливает всю систему: инфраструктуру и сервисы |
| `task restart-all` | Перезапускает всю систему: инфраструктуру и сервисы |

#### Управление сервисами
| Команда | Описание |
|---------|----------|
| `task start-backend` | Запускает только backend-api сервис |
| `task start-cron` | Запускает только cron-service сервис |
| `task start-email` | Запускает только email-sender сервис |
| `task start-all` | Запускает все сервисы проекта |
| `task stop-backend` | Останавливает backend-api сервис |
| `task stop-cron` | Останавливает cron-service сервис |
| `task stop-email` | Останавливает email-sender сервис |

#### Сборка образов
| Команда | Описание |
|---------|----------|
| `task build-backend` | Собирает Docker-образ для backend-api |
| `task build-cron` | Собирает Docker-образ для cron-service |
| `task build-email` | Собирает Docker-образ для email-sender |
| `task build-all` | Собирает Docker-образы для всех сервисов |

#### Логи и отладка
| Команда | Описание |
|---------|----------|
| `task logs-all` | Просмотр логов всех сервисов в реальном времени |
| `task logs-backend` | Просмотр логов backend-api |
| `task logs-cron` | Просмотр логов cron-service |
| `task logs-email` | Просмотр логов email-sender |

### Быстрый старт

1. Запустите инфраструктуру:
   ```bash
   task start-infra
   ```

2. Соберите сервисы:
   ```bash
   task build-all
   ```

3. Запустите все сервисы:
   ```bash
   task start-all
   ```

API будет доступен на `http://localhost:8080`

## Архитектура

Система построена на микросервисной архитектуре. Подробное описание и диаграмма доступны в папке [`architecture/`](architecture/):
- [`diagram.md`](architecture/diagram.md) - диаграмма в формате Mermaid
- [`diagram.png`](architecture/diagram.png) - визуальная диаграмма архитектуры

### Основные компоненты

- **Backend API** (`backend-api/`): Основной REST и gRPC API для управления пользователями и задачами
- **Cron Service** (`cron-service/`): Сервис для выполнения запланированных задач
- **Email Sender** (`email-sender/`): Сервис для асинхронной отправки email-уведомлений
- **API Proto** (`api-proto/`): Определения gRPC сервисов и протоколов

## Структура проекта

```
task-tracker/
├── backend-api/          # Основной API-сервис (HTTP/gRPC)
├── cron-service/         # Сервис для cron-задач
├── email-sender/         # Сервис отправки email
├── api-proto/            # gRPC определения и генерация кода
├── architecture/         # Диаграммы и документация архитектуры
├── docker-compose.yaml   # Конфигурация Docker Compose
├── Taskfile.yaml         # Автоматизация команд
├── .env.example          # Пример переменных окружения
└── README.md             # Этот файл
```

## Конфигурация

Все настройки производятся через переменные окружения. Скопируйте `.env.example` в `.env` и настройте:

- **Базы данных**: Параметры подключения к PostgreSQL (основная БД, БД для cron-сервиса, БД для email-sender)
- **API**: Порты для HTTP и gRPC серверов
- **JWT**: Секретный ключ для токенов аутентификации
- **Kafka**: Адреса брокеров и названия топиков
- **SMTP**: Настройки для отправки email

Подробное описание всех переменных см. в файле `.env.example`.
