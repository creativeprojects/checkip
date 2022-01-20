GOCMD=env go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get

BINARY=checkip
DEPLOY_SERVER=
DEPLOY_BINARY=/opt/checkip/checkip
BUILD_PROD=./build
TESTS=./...
COVERAGE_FILE=coverage.out

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: all test build coverage clean build-prod deploy

all: test build

build:
		$(GOBUILD) -o $(BINARY) -v

test:
		$(GOTEST) -v $(TESTS)

clean:
		$(GOCLEAN)
		rm -f $(BINARY) $(COVERAGE_FILE) ${BUILD_PROD}/$(BINARY)

build-prod:
		GOOS="linux" GOARCH="amd64" $(GOBUILD) -o ${BUILD_PROD}/$(BINARY) -v

deploy: build-prod
		rsync -avz ${BUILD_PROD}/$(BINARY) $(DEPLOY_SERVER):$(DEPLOY_BINARY)
