MUSL_LINK=https://musl.cc/x86_64-linux-musl-cross.tgz
MUSL_DIR=muslgcc
MUSL_BIN=$(MUSL_DIR)/bin/x86_64-linux-musl-gcc

CGO_ENABLED=1
GOOS=linux
GOARCH=amd64
CC="$(CURDIR)/$(MUSL_BIN)"
FLAG=-s -w --extldflags -static

all: build

deps:
	go mod tidy

$(MUSL_BIN):
	mkdir -p $(MUSL_DIR)
	wget -O $(MUSL_DIR).tgz $(MUSL_LINK)
	tar -zxf $(MUSL_DIR).tgz --strip-components=1 -C $(MUSL_DIR)

build: deps $(MUSL_BIN)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) CC=$(CC) go build -ldflags=$(FLAG)

clean:
	rm -rf $(MUSL_DIR) $(MUSL_DIR).tgz

.PHONY: all deps build clean
