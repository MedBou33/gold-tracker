.PHONY: fmt vet test build check

fmt:
	@test -z "$$(gofmt -l .)" || (echo "unformatted files:"; gofmt -l .; exit 1)

vet:
	go vet ./...

test:
	go test -v -cover ./...

build:
	go build -o bin/tracker ./cmd/tracker

check: fmt vet test build
