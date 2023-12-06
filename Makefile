# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=owl
ENV_VARS=DATABASE_URL=postgres://postgres:hoothoo@localhost:5432/owl_test?sslmode=disable

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	${ENV_VARS} $(GOTEST) -v ./... 

clean:
	rm -f $(BINARY_NAME)

.PHONY: all build test clean
