SHELL := /bin/sh

.PHONY: up down test lint seed

up:
	docker compose up --build

down:
	docker compose down -v

test:
	cd backend && go test ./...

lint:
	cd backend && go vet ./...

seed:
	@echo "No seed data yet. Placeholder for future scripts."
