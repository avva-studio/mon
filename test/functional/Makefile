REPO_NAME ?= accounting-rest

DOCKER_IMAGE ?= accounting-rest-serve:develop-v0.3.0-2-gc66e16c

TEST_BIN_NAME ?= functional.test

PROJECT_NAME ?= $(REPO_NAME)-functional-test

VERSION ?= $(shell git describe --tags --always)
LDFLAGS = "-w -X main.Version=$(VERSION)"

OS ?= linux
ARCH ?= amd64

all:
	$(MAKE) fetch-bin
	$(MAKE) build-tests
	$(MAKE) test; $(MAKE) clean

clean:
	rm -v ./bin/$(TEST_BIN_NAME)
	$(MAKE) docker-compose-down

build-tests:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go test -c -o ./bin/$(TEST_BIN_NAME) -installsuffix cgo -a -ldflags $(LDFLAGS)

test:
	$(MAKE) docker-compose-up

docker-compose-up:
	DOCKER_IMAGE=${DOCKER_IMAGE} \
	docker-compose --project-name $(PROJECT_NAME) up --build --abort-on-container-exit --force-recreate

docker-compose-down:
	DOCKER_IMAGE=${DOCKER_IMAGE} \
	docker-compose --project-name $(PROJECT_NAME) down