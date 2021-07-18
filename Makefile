all: binaries

binaries:
	env CGO_ENABLED=0 go build -o ./bin/ ./cmd/...

test: lint vet
	go test -race -cover ./...

test_bench: lint vet
	go test -race -cover -bench=. ./...

lint:
	staticcheck ./...

vet:
	go vet -v ./...

install:
	go install -v ./cmd/...

generate:
	go generate -v ./...

.PHONY: all binaries test test_bench lint vet install generate
