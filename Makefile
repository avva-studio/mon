VERSION ?= $(shell git describe --tags --always)
LDFLAGS = "-w -X main.Version=$(VERSION)"
NAME = accounting-rest-serve

OS ?= linux
ARCH ?= amd64

build-bin:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(NAME) -a -ldflags $(LDFLAGS) ./cmd/serve

build-docker:
	docker build --tag $(NAME):$(VERSION) .
	