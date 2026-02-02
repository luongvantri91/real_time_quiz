.PHONY: help deps build run test clean docker-up docker-down docker-logs

help:
	@echo "RT-Quiz Makefile Commands"
	@echo "=========================="
	@echo "make deps        - Download Go dependencies"
	@echo "make build       - Build binary"
	@echo "make run         - Run server locally"
	@echo "make test        - Run tests"
	@echo "make clean       - Remove binary and build artifacts"
	@echo "make docker-up   - Start Docker containers (Redis + Postgres)"
	@echo "make docker-down - Stop Docker containers"
	@echo "make docker-logs - View Docker logs"
	@echo "make docker-build - Build Docker image"

deps:
	go mod download
	go mod tidy

build: deps
	go build -o rt-quiz .

run: deps
	go run main.go

test:
	go test ./... -v -race

clean:
	rm -f rt-quiz rt-quiz.exe
	go clean

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-build:
	docker build -t rt-quiz:latest .

docker-run: docker-build
	docker run --network host \
		-e REDIS_URL=redis://localhost:6379/0 \
		-e POSTGRES_URL=postgres://quiz_user:quiz_password@localhost:5432/rt_quiz \
		rt-quiz:latest
