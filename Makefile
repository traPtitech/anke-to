.PHONY: dev
dev:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	docker-compose -f development/docker-compose.yaml up --build