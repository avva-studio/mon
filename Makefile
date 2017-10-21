VERSION ?= $(shell git describe --tags --always)
LDFLAGS = "-w -X main.Version=$(VERSION)"
NAME = gohmoney-rest

OS ?= linux
ARCH ?= amd64

build:
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(NAME) -a -ldflags $(LDFLAGS)

test:
	go list ./... |grep -v vendor | xargs go test -v

docker-compose:
	@$(MAKE) docker-compose-build
	@$(MAKE) docker-compose-up

docker-compose-build:
	@$(MAKE) build
	@VERSION=$(VERSION) docker-compose -f docker-compose.yml build

docker-compose-up:
	@VERSION=$(VERSION) docker-compose -f docker-compose.yml up -d

coverage:
	go test -coverprofile=coverage.out && go tool cover -func=coverage.out && go tool cover -html=coverage.out
