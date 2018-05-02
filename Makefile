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
	$(MAKE) build-monserve
	$(MAKE) build-monserve-docker

build-monserve:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(SERVE_NAME) -a -ldflags $(LDFLAGS) ./cmd/$(SERVE_NAME)

build-monserve-docker:
	docker build --tag $(SERVE_NAME):$(VERSION) .

moncli:
	$(MAKE) build-moncli

build-moncli:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(CLI_NAME) -a -ldflags $(LDFLAGS) ./cmd/$(CLI_NAME)
