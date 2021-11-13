KAFKA_GROUP_ID ?= $(shell echo "otel-example-$$(date +%s%3)")
KAFKA_BOOTSTRAP_SERVERS ?=http://host.docker.internal:9092
KAFKA_TOPIC ?= otel-example

GOLINTER ?= revive

.PHONY: deps
deps:
	go mod download

.PHONY: run-root-http
run-root: deps
	go run cmd/root/main.go

.PHONY: run-http
run-http: deps
	PORT=1618 go run cmd/http/main.go

.PHONY: run-events
run-events: deps
	go run cmd/event/main.go

.PHONY: lint
lint:
	${GOLINTER} ./...

.PHONY: compose-dependencies
compose-dependencies:
	docker-compose up \
    		--remove-orphans \
    		-d \
    		elasticsearch kafka zookeeper jaeger

.PHONY: compose-down
compose-down:
	docker-compose down

.PHONY: clean
clean:
	docker-compose down --volumes --remove-orphans
