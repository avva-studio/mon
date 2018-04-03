VERSION ?= $(shell git describe --tags --always)
LDFLAGS = "-w -X main.Version=$(VERSION)"

REPO_NAME = accounting-rest
SERVE_NAME = $(REPO_NAME)-serve
CLI_NAME = $(REPO_NAME)-cli

OS ?= linux
ARCH ?= amd64

build:
	$(MAKE) build-serve
	$(MAKE) build-cli
	$(MAKE) build-serve-docker

build-serve:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(SERVE_NAME) -a -ldflags $(LDFLAGS) ./cmd/serve

build-serve-docker:
	docker build --tag $(SERVE_NAME):$(VERSION) .

build-cli:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(CLI_NAME) -a -ldflags $(LDFLAGS) ./cmd/cli