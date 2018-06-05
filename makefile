
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GODEP=dep
BINARY_NAME=hugo-updater
BINARY_LINUX=$(BINARY_NAME)-linux-amd64

all: test build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	$(GODEP) ensure



# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /Users/ndoye/go/src/hugo-updater golang:latest go build -o "$(BINARY_LINUX)" -v    
