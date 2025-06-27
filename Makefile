CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
FLAG="-s -w"

all: build

deps:
	go mod tidy

build: deps
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(FLAG)

.PHONY: all deps build
