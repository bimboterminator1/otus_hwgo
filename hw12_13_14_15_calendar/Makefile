BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... 

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.3

lint: install-lint-deps
	golangci-lint -c ./.golangci.yml run ./...

POSTGRES_USER := postgres
POSTGRES_PASSWORD := postgres
POSTGRES_DB := calendar
APP_USER := appuser
APP_PASSWORD := apppass
POSTGRES_PORT := 5432
CONTAINER_NAME := calendar_postgres

postgres:
	docker run --name $(CONTAINER_NAME) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-p $(POSTGRES_PORT):5432 \
		-d postgres && \
	sleep 3 && \
	docker exec -it $(CONTAINER_NAME) psql -U $(POSTGRES_USER) -c "CREATE USER $(APP_USER) WITH PASSWORD '$(APP_PASSWORD)';" && \
	docker exec -it $(CONTAINER_NAME) psql -U $(POSTGRES_USER) -c "CREATE DATABASE $(POSTGRES_DB);" && \
	docker exec -it $(CONTAINER_NAME) psql -U $(POSTGRES_USER) -c "GRANT ALL PRIVILEGES ON DATABASE $(POSTGRES_DB) TO $(APP_USER);" && \
	docker exec -it $(CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "GRANT ALL ON SCHEMA public TO $(APP_USER);"

pg-down:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

generate:
	protoc -I ./api/proto ./api/proto/*.proto --go_out=./api/gen/ \
	 --go-grpc_out=./api/gen/ 
	

.PHONY: build run build-img run-img version test lint postgres pg-down
