test:
	go test ./...

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all