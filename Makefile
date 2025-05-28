# Galaxy Node Pool Makefile
# AI-ID: CP-GAL-NODEPOOL-001

# Variables
BINARY_NAME=galaxy-pool
VERSION=$(shell cat VERSION)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitCommit=$(GIT_COMMIT) -X main.GitBranch=$(GIT_BRANCH)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build
.PHONY: build
build:
	mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

# Install
.PHONY: install
install:
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

# Build for Linux
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-linux ./cmd/$(BINARY_NAME)

# Build for Windows
.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-windows.exe ./cmd/$(BINARY_NAME)

# Build for MacOS
.PHONY: build-darwin
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin ./cmd/$(BINARY_NAME)

# Build all platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin

# Clean
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf bin/

# Test
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run
.PHONY: run
run:
	$(GOBUILD) -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	./bin/$(BINARY_NAME)

# Get dependencies
.PHONY: deps
deps:
	$(GOGET) -u ./...

# Update dependencies
.PHONY: update-deps
update-deps:
	$(GOMOD) tidy

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Create release
.PHONY: release
release: build-all
	@echo "Creating release $(VERSION)"
	tar -czf bin/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C bin $(BINARY_NAME)-linux
	zip -j bin/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip bin/$(BINARY_NAME)-windows.exe
	tar -czf bin/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz -C bin $(BINARY_NAME)-darwin

# Show help
.PHONY: help
help:
	@echo 'Available targets:'
	@echo '  build         - Build the binary for current platform'
	@echo '  install       - Install the binary to GOPATH/bin'
	@echo '  build-linux   - Build for Linux'
	@echo '  build-windows - Build for Windows'
	@echo '  build-darwin  - Build for MacOS'
	@echo '  build-all     - Build for all platforms'
	@echo '  clean         - Remove binary and build artifacts'
	@echo '  test          - Run tests'
	@echo '  run           - Build and run the application'
	@echo '  deps          - Get dependencies'
	@echo '  update-deps   - Update dependencies'
	@echo '  lint          - Run linter'
	@echo '  release       - Create release packages'
	@echo '  help          - Show this help message'
