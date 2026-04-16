.PHONY: init
init:
	go mod download
	go install go.uber.org/mock/mockgen
	go install github.com/google/wire/cmd/wire

generate:

.PHONY: dev
dev:
	docker-compose -f docker/dev/docker-compose.yaml up --build

.PHONY: test
test:
	-docker-compose -f docker/test/docker-compose.yaml down
	docker-compose -f docker/test/docker-compose.yaml up --build

.PHONY: build
build:
	go build -o anke-to