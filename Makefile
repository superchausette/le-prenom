build-server:
	go build -o bin/server cmd/server/main.go


build-populate-db:
	go build -o bin/populate_bd cmd/populate_db/main.go

build: build-server build-populate-db

test-race:
	go test -race ./...

test-coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...

test: test-race test-coverage

dep:
	go mod download

vet:
	go vet

format:
	goimports -w $(shell find . -type f -name '*.go')

lint:
	golangci-lint run --enable-all