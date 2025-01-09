.PHONY: all build run test clean

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOTEST = $(GOCMD) test
GOCLEAN = $(GOCMD) clean
BINARY_NAME = cstv-go

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/main.go

run:
	$(GORUN) ./cmd/main.go

test:
	$(GOTEST) ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
