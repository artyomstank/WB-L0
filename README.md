# Order Service (WB Tech Demo)

Микросервис для просмотра данных о заказе с использованием Go, PostgreSQL, Kafka и кэширования.

## Установка и запуск

### Предварительные требования

* Go 1.19+
* Docker и Docker Compose
* PostgreSQL 14+
* Kafka


## Быстрый старт

### Запуск всей инфраструктуры и сервиса:
**Запуск**
```bash
make build
make up
make create-topic
```
**Полная перезагрузка**
   ```bash
   make down
   make clean
   ```

## API Endpoints

### GET /api/order/{uid}
Получение заказа по ID.

**Пример запроса:**
```bash
curl -X GET http://localhost:8081/api/order/b563feb7b2b84b6test
```

**Успешный ответ (200 OK):**
```json
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```



### GET /health
Проверка работоспособности сервиса.

**Пример запроса:**
```bash
curl -X GET http://localhost:8081/health
```

**Успешный ответ (200 OK):**
```json
{
  "status": "OK",
  "timestamp": "2023-10-01T12:00:00Z"
}
```



## Переменные окружения

Перед запуском убедитесь, что установлены необходимые переменные окружения. Основные переменные:

- `KAFKA_HOST` - хост Kafka (по умолчанию: localhost)
- `KAFKA_PORT` - порт Kafka (по умолчанию: 9092)
- `KAFKA_TOPIC` - топик Kafka (по умолчанию: wb-orders)
- `KAFKA_GROUP` - группа Kafka (по умолчанию: wb-tech-demo-service)

## Управление через Makefile

### Docker Compose операции
```bash
make up          # Запуск всех сервисов
make down        # Остановка всех сервисов
make build       # Пересборка образов
make logs        # Просмотр логов
make ps          # Статус контейнеров
make restart     # Перезапуск сервисов
make clean       # Полная очистка (с удалением volumes)
```

### Kafka операции
```bash
make create-topic     # Создание топика
make list-topics      # Список топиков
make describe-topic   # Описание топика
make delete-topic     # Удаление топика
```

### Сервисные команды
```bash
make start-service    # Запуск только основного сервиса
make stop-service     # Остановка сервиса
make restart-service  # Перезапуск сервиса
make start-producer   # Запуск продюсера
make stop-producer    # Остановка продюсера
```

### База данных
```bash
make list-tables      # Просмотр таблиц БД
make view-f5          # Просмотр первых 5 записей (указать TABLE=table_name)
```

### Тесты
```bash
make test             # Запуск всех тестов
make test-coverage    # Запуск тестов с покрытием
make generate-mocks   # Генерация моков для тестов
```
## Подробнее про команды make
### Установка зависимостей
```bash
make deps
```

### Тестирование

#### Запуск всех тестов:
```bash
make test
```

#### Проверка покрытия тестами:
```bash
make test-coverage
```

#### Генерация моков для тестов:
```bash
make generate-mocks
```

### Запуск проекта

#### Запуск всей инфраструктуры (Docker Compose):
```bash
make up
```

#### Сборка и пересборка образов:
```bash
make build
```

#### Запуск только основного сервиса:
```bash
make start-service
```

#### Запуск продюсера:
```bash
make start-producer
```

#### Остановка сервиса:
```bash
make stop-service
```

#### Перезапуск сервиса:
```bash
make restart-service
```

### Управление Kafka

#### Создание топика:
```bash
make create-topic
```

#### Просмотр списка топиков:
```bash
make list-topics
```

#### Описание топика:
```bash
make describe-topic
```

#### Удаление топика:
```bash
make delete-topic
```

### Работа с базой данных

#### Просмотр таблиц БД:
```bash
make list-tables
```

#### Просмотр первых 5 записей таблицы:
```bash
make view-f5 TABLE=orders
```

### Мониторинг и логи

#### Просмотр логов всех сервисов:
```bash
make logs
```

#### Статус контейнеров:
```bash
make ps
```

### Администрирование

#### Полная остановка и очистка:
```bash
make down
```

#### Полная очистка с удалением volumes:
```bash
make clean
```

#### Перезапуск всех сервисов:
```bash
make restart
```

## Тестирование
#### Проект включает следующие виды тестов:
```bash
# Конкретный тест
go test -v -run TestGetOrderByUID

# Тесты с таймаутом
go test -v -timeout 30s ./...

# Параллельные тесты
go test -v -parallel 4 ./...
```

## Архитектура
```
                                            +---cmd
                                            |   |   main.go
                                            |   |
                                            |   \---producer
                                            |           st_prod.go
                                            |
                                            +---config
                                            |       config.go
                                            |       config_test.go
                                            |
                                            +---internal
                                            |   +---cache
                                            |   |       cache.go
                                            |   |       cache_test.go
                                            |   |
                                            |   +---db
                                            |   |       db.go
                                            |   |
                                            |   +---generator
                                            |   |       order.go
                                            |   |       order_test.go
                                            |   |
                                            |   +---handler
                                            |   |       handler.go
                                            |   |       handler_interface.go
                                            |   |       handler_test.go
                                            |   |       server.go
                                            |   |
                                            |   +---kafka
                                            |   |       consumer.go
                                            |   |       consumer_test.go
                                            |   |       interfaces.go
                                            |   |       producer.go
                                            |   |       producer_test.go
                                            |   |
                                            |   +---mocks
                                            |   |       mock_cache.go
                                            |   |       mock_repository.go
                                            |   |       mock_service.go
                                            |   |
                                            |   +---models
                                            |   |       models_test.go
                                            |   |       order.go
                                            |   |       validator.go
                                            |   |
                                            |   +---repo
                                            |   |       postgres.go
                                            |   |       repository.go
                                            |   |       repository_interface.go
                                            |   |       repository_test.go
                                            |   |       transaction_test.go
                                            |   |
                                            |   \---service
                                            |           service.go
                                            |           service_interface.go
                                            |           service_test.go
                                            |
                                            +---migrations
                                            |       20250829195353_add_orders_table.down.sql
                                            |       20250829195353_add_orders_table.up.sql
                                            |
                                            \---web
                                                    index.html
                                            .env
                                            .gitignore
                                            coverage.out
                                            docker-compose.yml
                                            Dockerfile
                                            Dockerfile.producer
                                            go.mod
                                            go.sum
                                            Makefile
                                            README.md
```

## Troubleshooting

### Частые проблемы
1. **Kafka недоступна** - проверьте `make ps` и `make create-topic`
2. **Миграции не применяются** - проверьте строку подключения к БД
3. **Кэш не заполняется** - проверьте настройки CACHE_STARTUP_SIZE

### Полезные команды для отладки
```bash
# Проверка подключения к Kafka
make list-topics

# Проверка данных в БД
make list-tables
make view-f5 TABLE=orders

# Перезапуск сервиса
make restart-service
```