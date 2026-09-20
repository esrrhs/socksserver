BINARY_NAME=socksserver
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "0.3")
BUILD_DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-s -w -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)

.PHONY: all build clean test race vet fmt docker pack pack-fast tidy

all: build

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) .

test:
	go test -v ./...

race:
	go test -v -race ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

docker:
	docker build -t $(BINARY_NAME):latest .

pack:
	./pack.sh

pack-fast:
	./pack.sh fast

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	rm -rf pack pack.zip
	go clean -testcache
