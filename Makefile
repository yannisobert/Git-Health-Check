.PHONY: dev dev-go dev-front build build-front test lint tidy

dev-go:
	go run ./cmd/owlspector server

dev-front:
	cd frontend && npm run dev

dev:
	make -j2 dev-go dev-front

build-front:
	cd frontend && npm run build
	rm -rf server/dist
	cp -r frontend/dist server/dist

build: build-front
	go build -o bin/owlspector ./cmd/owlspector

test:
	go test ./cmd/... ./internal/... ./server/...

lint:
	golangci-lint run

tidy:
	go mod tidy
