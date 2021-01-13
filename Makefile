.PHONY: dev
dev:
	docker-compose -f docker/dev/docker-compose.yaml up --build

.PHONY: test
test:
	-docker-compose -f docker/test/docker-compose.yaml down
	docker-compose -f docker/test/docker-compose.yaml up --build

.PHONY: tuning
tuning:
	docker-compose -f docker/tuning/docker-compose.yaml up --build

.PHONY: pprof
pprof:
	go tool pprof -png -output pprof.png http://localhost:6060/debug/pprof/profile

.PHONY: slow
slow:
	docker-compose -f docker/tuning/docker-compose.yaml exec mysql pt-query-digest /tmp/mysql-slow.sql

.PHONY: myprof
myprof:
	docker-compose -f docker/tuning/docker-compose.yaml exec mysql myprofiler -user=root -password=password ${ARGS}

.PHONY: build
build:
	go build -o anke-to

.PHONY: bench
bench: build
	./anke-to bench

.PHONY: bench-init
bench-init: build
	./anke-to init