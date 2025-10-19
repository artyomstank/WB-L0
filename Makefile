DOCKER_COMPOSE = docker compose
KAFKA_CONTAINER = kafka
BROKER = $(KAFKA_HOST):$(KAFKA_PORT)

# Default values (can be overridden via env)
KAFKA_HOST ?= localhost
KAFKA_PORT ?= 9092
KAFKA_TOPIC ?= wb-orders
KAFKA_GROUP ?= wb-tech-demo-service

.PHONY: up down build logs ps create-topic list-topics describe-topic delete-topic restart clean test test-coverage generate-mocks

# Docker compose commands
up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

build:
	$(DOCKER_COMPOSE) build

logs:
	$(DOCKER_COMPOSE) logs -f

ps:
	$(DOCKER_COMPOSE) ps

restart:
	$(DOCKER_COMPOSE) restart

clean:
	$(DOCKER_COMPOSE) down -v --remove-orphans

# Kafka commands
create-topic:
	$(DOCKER_COMPOSE) exec -T $(KAFKA_CONTAINER) \
	kafka-topics.sh --create --if-not-exists \
	--bootstrap-server $(BROKER) \
	--replication-factor 1 \
	--partitions 3 \
	--topic $(KAFKA_TOPIC)

list-topics:
	$(DOCKER_COMPOSE) exec -T $(KAFKA_CONTAINER) \
	kafka-topics.sh --list --bootstrap-server $(BROKER)

describe-topic:
	$(DOCKER_COMPOSE) exec -T $(KAFKA_CONTAINER) \
	kafka-topics.sh --describe --bootstrap-server $(BROKER) --topic $(KAFKA_TOPIC)

delete-topic:
	$(DOCKER_COMPOSE) exec -T $(KAFKA_CONTAINER) \
	kafka-topics.sh --delete --bootstrap-server $(BROKER) --topic $(KAFKA_TOPIC)

# Service commands
start-service:
	$(DOCKER_COMPOSE) up -d wb-service

stop-service:
	$(DOCKER_COMPOSE) stop wb-service

restart-service:
	$(DOCKER_COMPOSE) restart wb-service

start-producer:
	$(DOCKER_COMPOSE) up -d wb-producer

stop-producer:
	$(DOCKER_COMPOSE) stop wb-producer

#Postgres
list-tables:
	$(DOCKER_COMPOSE) exec postgres psql -U wb_user -d wb_demo_db -c "\dt"
view-f5:
	$(DOCKER_COMPOSE) exec postgres psql -U wb_user -d wb_demo_db -c "SELECT * FROM $(TABLE) LIMIT 5;"

# Testing
.PHONY: test test-coverage generate-mocks

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

generate-mocks:
	mockgen -source=internal/repo/repository_interface.go -destination=internal/mocks/mock_repository.go -package=mocks
	mockgen -source=internal/service/service_interface.go -destination=internal/mocks/mock_service.go -package=mocks
	mockgen -source=internal/cache/cache.go -destination=internal/mocks/mock_cache.go -package=mocks
