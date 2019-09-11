.PHONY: dev
dev:
	go mod download
	docker-compose -f development/docker-compose.yaml up --build
