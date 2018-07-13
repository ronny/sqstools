test: vet fmt test-only

install:
	go install -v ./cmd/...

build:
	gox -os="linux darwin windows" -arch="amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" github.com/ronny/sqstools/cmd/...
	upx bin/*

vet:
	go vet -v ./...

fmt:
	go fmt ./...

test-only:
	go test -v -race -cover -bench=. ./...

mocks:
	mockery -all -dir ./internal/sqstools/ -output ./internal/mocks/

.PHONY: test install build vet fmt test-only mocks
