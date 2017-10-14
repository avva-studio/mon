VERSION ?= $(shell git describe --tags --always)

IMAGE = "pocketmedia/offer-api"
PKG = "github.com/Pocketbrain/offer-api"

LDFLAGS = "-w -X main.Version=$(VERSION)"

NAME = gohmoney-rest

OS ?= linux
ARCH ?= amd64

build:
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -installsuffix cgo -o bin/$(NAME) -a -ldflags $(LDFLAGS)

test:
	go list ./... |grep -v vendor | xargs go test -v

docker-compose:
	@$(MAKE) build
	docker-compose -f docker-compose.yml up --build

coverage:
	go test -coverprofile=coverage.out && go tool cover -func=coverage.out && go tool cover -html=coverage.out

# serve-local:
	# go run main.go serve -H 127.0.0.1 -p 8083