DOCKER_COMPOSE = docker-compose
KAFKA_CONTAINER = kafka
BROKER = $(KAFKA_HOST):$(KAFKA_PORT)

# Значения по умолчанию (можно переопределять через env)
KAFKA_HOST ?= localhost
KAFKA_PORT ?= 9092
KAFKA_TOPIC ?= wb-orders
KAFKA_GROUP ?= wb-tech-demo-service

.PHONY: create-topic list-topics describe-topic delete-topic

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
