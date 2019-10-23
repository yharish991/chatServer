# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOLINT=$(GOCMD)lint
GOTEST=$(GOCMD) test
BINARY_NAME=chatserverbinary
GOPACKAGES=$(shell go list ./...)

all: deps lint test run
deps:
	go get -u github.com/onsi/ginkgo/ginkgo  # installs the ginkgo CLI
	go get -u github.com/onsi/gomega/...     # fetches the matcher library
	go get github.com/stretchr/testify
lint:
	$(GOLINT) --set_exit_status ${GOPACKAGES}
test:
	$(GOTEST) -v ./...
build:
	$(GOBUILD) -o $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME)
	./$(BINARY_NAME)