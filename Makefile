GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=core

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) cmd/cli/main.go

run:
	$(MAKE) build
	./$(BINARY_NAME)