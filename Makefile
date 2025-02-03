
build:
	go build -o bin/main ./cmd

test:
	go test ./...

docker-build:
	docker-compose build

run:
	 docker-compose up
