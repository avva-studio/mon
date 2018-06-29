VERSION ?= $(shell git describe --tags --always)
LDFLAGS = "-w -X main.Version=$(VERSION)"

SERVE_NAME = monserve
CLI_NAME = moncli

OS ?= linux
ARCH ?= amd64

all: build install clean

build:
	$(MAKE) monserve
	$(MAKE) moncli

install:
	cp ./bin/* $(GOPATH)/bin/

clean:
	rm ./bin/*

monserve:
	$(MAKE) monserve-binary
	$(MAKE) monserve-image

monserve-binary:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(SERVE_NAME) -a -ldflags $(LDFLAGS) ./cmd/$(SERVE_NAME)

monserve-image:
	docker build --tag $(SERVE_NAME):$(VERSION) .

moncli:
	$(MAKE) build-moncli

build-moncli:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(CLI_NAME) -a -ldflags $(LDFLAGS) ./cmd/$(CLI_NAME)
