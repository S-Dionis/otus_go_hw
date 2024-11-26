BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

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
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc --proto_path=api/calendar/ --go_out=pb --go-grpc_out=pb api/*.proto

up:
	docker compose -f ./deployments/docker-compose.yaml up -d

down:
	docker compose down

integration-tests: up run-tests down

run-tests:
	docker exec -it test go test /go/src/test/test -v
	@if [ $$? -eq 0 ]; then \
		echo "Integration tests passed."; \
		exit 0; \
	else \
		echo "Integration tests failed."; \
		exit 1; \
	fi

.PHONY: build run build-img run-img version test lint

