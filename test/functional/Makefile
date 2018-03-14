REPO_NAME ?= accounting-rest
REPO_URI ?= github.com/glynternet/$(REPO_NAME)

BIN_NAME ?= $(REPO_NAME)

PROJECT_NAME ?= $(REPO_NAME)-functional-test

all:
	$(MAKE) fetch-bin
	$(MAKE) docker-compose-build-up

fetch-bin:
	cp -v $(GOPATH)/src/$(REPO_URI)/bin/$(BIN_NAME) ./bin/$(BIN_NAME)

docker-compose-build-up:
	docker-compose --project-name $(PROJECT_NAME) up --build