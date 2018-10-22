GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
BINARY_NAME=core

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) cmd/cli/main.go

tidy:
	$(GOMOD) tidy

run:
	$(MAKE) build
	./$(BINARY_NAME)