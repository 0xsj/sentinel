.PHONY: install dev build check frontend

install:
	npm --prefix frontend ci
	go mod download

dev:
	go tool wails dev

build:
	go tool wails build

check:
	npm --prefix frontend run check
	npm --prefix frontend run build
	go vet ./...

frontend:
	npm --prefix frontend run dev
